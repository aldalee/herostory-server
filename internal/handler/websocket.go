package handler

import (
	"encoding/binary"
	"errors"
	"herostory-server/internal/pb"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	ErrUpgradeWebSocket = errors.New("websocket upgrade failed")
	ErrReadMessage      = errors.New("websocket read message failed")
)

func WebSocketHandshake(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error().Msgf("%v: %v", ErrUpgradeWebSocket, err)
		return
	}
	defer conn.Close()

	log.Info().Msgf("client %v connected to websocket", conn.RemoteAddr())

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Error().Msgf("%v: %v", ErrReadMessage, err)
			break
		}

		log.Info().Msgf("received client %v message: %v", conn.RemoteAddr(), data)

		code := binary.BigEndian.Uint16(data[2:4])
		msg, err := pb.DecodeMessage(data[4:], int16(code))
		if err != nil {
			log.Error().Msgf("decode client %v message failed, code: %v, err: %v", conn.RemoteAddr(), code, err)
			continue
		}

		log.Info().Msgf("decode client %v message success, code: %v, msg: %v", conn.RemoteAddr(), code, msg.Descriptor().Name())

		desc := pb.File_api_proto_game_msg_proto.Messages().ByName("UserLoginCmd")
		username := msg.Get(desc.Fields().ByName("userName"))
		password := msg.Get(desc.Fields().ByName("password"))
		log.Info().Msgf("client %v login with username: %v, password: %v", conn.RemoteAddr(), username, password)

		rest := &pb.UserLoginResult{
			UserId:     1,
			UserName:   username.String(),
			HeroAvatar: "Hero_Shaman",
		}
		byteArray, err := pb.EncodeMessage(rest)
		if err != nil {
			log.Error().Msgf("encode client %v login result failed, err: %v", conn.RemoteAddr(), err)
			continue
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, byteArray); err != nil {
			log.Error().Msgf("write client %v login result failed, err: %v", conn.RemoteAddr(), err)
		}
	}
}
