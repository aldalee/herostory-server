package handler

import (
	"errors"
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
	}
}
