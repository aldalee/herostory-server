package server

import (
	"context"
	"encoding/binary"
	"herostory-server/internal/codec"
	"herostory-server/internal/pb"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

const wsAddr = "ws://localhost:12345/websocket"

func TestMain(m *testing.M) {
	codec.InitMaps()
	os.Exit(m.Run())
}

func dial(t *testing.T) *websocket.Conn {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()
	conn, resp, err := (&websocket.Dialer{}).DialContext(ctx, wsAddr, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
	t.Cleanup(func() { conn.Close() })
	return conn
}

func send(t *testing.T, conn *websocket.Conn, msg proto.Message) {
	t.Helper()
	data, err := codec.EncodeMessage(msg)
	require.NoError(t, err)
	require.NoError(t, conn.WriteMessage(websocket.BinaryMessage, data))
}

func recv(t *testing.T, conn *websocket.Conn) (uint16, []byte) {
	t.Helper()
	deadline, _ := t.Context().Deadline()
	if deadline.IsZero() {
		deadline = time.Now().Add(10 * time.Second)
	}
	_ = conn.SetReadDeadline(deadline)
	_, raw, err := conn.ReadMessage()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(raw), 4)
	return binary.BigEndian.Uint16(raw[2:4]), raw[4:]
}

func recvAs[T any, PT interface {
	*T
	proto.Message
}](t *testing.T, conn *websocket.Conn, wantCode pb.MsgCode) *T {
	t.Helper()
	code, body := recv(t, conn)
	assert.Equal(t, uint16(wantCode), code)
	msg := PT(new(T))
	require.NoError(t, proto.Unmarshal(body, msg))
	return (*T)(msg)
}

func TestWebSocketConnection(t *testing.T) {
	conn := dial(t)
	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte("ping")))
}

func TestParseUserLoginCmd(t *testing.T) {
	want := &pb.UserLoginCmd{UserName: "test_player", Password: "test_pass123"}
	data, err := proto.Marshal(want)
	require.NoError(t, err)

	got := new(pb.UserLoginCmd)
	require.NoError(t, proto.Unmarshal(data, got))
	assert.True(t, proto.Equal(want, got))
}

func TestUserLogin(t *testing.T) {
	for _, tt := range []struct {
		name, user, pass string
	}{
		{"existing_user", "EEE", "123456"},
		{"new_user", "Alice", "abc123"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			conn := dial(t)
			send(t, conn, &pb.UserLoginCmd{UserName: tt.user, Password: tt.pass})

			result := recvAs[pb.UserLoginResult](t, conn, pb.MsgCode_USER_LOGIN_RESULT)
			t.Logf("userId=%d userName=%s heroAvatar=%s", result.UserId, result.UserName, result.HeroAvatar)

			assert.NotZero(t, result.UserId)
			assert.Equal(t, tt.user, result.UserName)
		})
	}
}
