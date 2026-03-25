package main

import (
	"herostory-server/internal/app"
	"herostory-server/internal/server"
	"net/http"
)

func main() {
	app.Init()
	http.HandleFunc("/health", server.HealthCheck)
	http.HandleFunc("/websocket", server.WebSocketHandshake)
	http.ListenAndServe(":12345", nil)
}
