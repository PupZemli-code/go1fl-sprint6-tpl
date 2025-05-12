package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

// Структура сервера
type Server struct {
	Logger     *log.Logger
	HTTPServer *http.Server
}

// Функция создания сервера
func NewServer(logger *log.Logger) *Server {
	// Создаем роутер
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(handlers.HomeHandler))
	router.Handle("/upload", http.HandlerFunc(handlers.UploadHandler))

	// Настраиваем сервер
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		Logger:     logger,
		HTTPServer: server, // Теперь храним указатель на сервер
	}
}

func InitLogger() *log.Logger {
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
