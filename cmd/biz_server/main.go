package main

import (
	"herostory-server/bootstrap"
	"herostory-server/internal/handler"
	"net/http"
)

func main() {
	bootstrap.InitApp()

	http.HandleFunc("/health", handler.HealthCheck)
	http.ListenAndServe(":12345", nil)
}
