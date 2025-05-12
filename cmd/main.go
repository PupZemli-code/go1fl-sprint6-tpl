package main

import (
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	// Создаем логгер
	logger := server.InitLogger()

	// Создаем сервер
	srv := server.NewServer(logger)

	// Запускаем сервер
	err := srv.HTTPServer.ListenAndServe()
	if err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
