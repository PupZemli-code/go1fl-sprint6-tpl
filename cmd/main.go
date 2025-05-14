package main

import (
	"io"
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

/*
	func initLogger() *log.Logger {
		// Создаем файл для логов
		logfile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("ошибка создания или открытия файла app.log: %v", err)
		}

		// Настраиваем форматирование логов
		logger := log.New(logfile, "LOG: ", log.Ldate|log.Ltime|log.Lmicroseconds) // Добавляем дату и время
		io.MultiWriter(logfile, os.Stdout)                                         // Пишем в файл и консоль
		defer logfile.Close()
		return logger
	}
*/
func main() {
	logfile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("ошибка создания или открытия файла app.log: %v", err)
	}
	// Настраиваем форматирование логов
	logger := log.New(logfile, "LOG: ", log.Ldate|log.Ltime|log.Lmicroseconds) // Добавляем дату и время
	io.MultiWriter(logfile, os.Stdout)                                         // Пишем в файл и консоль
	defer logfile.Close()

	// Создаем сервер
	srv := server.NewServer(logger)
	logger.Println("Сервер Создан")
	// Запускаем сервер
	err = srv.HTTPServer.ListenAndServe()
	if err != nil {
		logger.Fatal("Ошибка запуска сервера:", err)
	}
	logger.Println("Сервер запущен")
}
