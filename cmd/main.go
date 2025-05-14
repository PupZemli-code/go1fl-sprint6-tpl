package main

import (
	"io"
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func initLogger() *log.Logger {
	// Создаем файл для логов
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Настраиваем форматирование логов
	logger := log.New(file, "LOG: ", log.Ldate|log.Ltime|log.Lmicroseconds) // Добавляем дату и время
	io.MultiWriter(file, os.Stdout)                                         // Пишем в файл и консоль

	return logger
}

func main() {
	// Создаем логгер
	logger := initLogger()

	// Создаем сервер
	srv := server.NewServer(logger)

	// Запускаем сервер
	err := srv.HTTPServer.ListenAndServe()
	if err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
