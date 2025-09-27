package handler

import (
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSocketConnection(t *testing.T) {
	dialer := websocket.Dialer{
		HandshakeTimeout: 3 * time.Second,
	}

	conn, resp, err := dialer.Dial("ws://localhost:12345/websocket", nil)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Close()

	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte("Hello World")))
}
