package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

// HomeHandler обрабатывает корневой путь и отдает индексный файл
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// Handler для эндпоинта /upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {

	logfile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("ошибка создания или открытия файла app.log: %v", err)
	}

	// Настраиваем форматирование логов
	logger := log.New(logfile, "LOG: ", log.Ldate|log.Ltime|log.Lmicroseconds) // Добавляем дату и время
	io.MultiWriter(logfile, os.Stdout)                                         // Пишем в файл и консоль
	defer logfile.Close()

	// 1. Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		logger.Println("Метод не поддерживается")
		return
	}
	logger.Println("Метод проверен успешно")

	// 2. Парсинг формы
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB ограничение
		http.Error(w, "Ошибка парсинга формы", http.StatusInternalServerError)
		logger.Println("Ошибка парсинга формы")
		return
	}
	logger.Println("Парсинг прошел успешно")

	// 3. Получение файла из формы
	fileHtml, fileHeader, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusInternalServerError)
		logger.Println("Файл не найден")
		return
	}
	defer fileHtml.Close()
	logger.Println("Файл из формы получен успешно")

	// 4. Чтение данных из файла
	data, err := io.ReadAll(fileHtml)
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		logger.Println("Ошибка чтения файла")
		return
	}
	logger.Println("Чтение данных из файла успешно")

	// 5. Автоопределение и конвертация
	convertedString, err := service.Service(string(data))
	if err != nil {
		http.Error(w, "Ошибка конвертации", http.StatusInternalServerError)
		logger.Println("Ошибка конвертации")
		return
	}
	logger.Println("Конвертация прошла успешно")

	// 6. Генерация имени файла, создание и запись в локальный файл
	fileName := time.Now().UTC().Format("2006-01-02_15-04-05")
	extension := filepath.Ext(fileHeader.Filename)
	// Проверяем, что расширение не пустое
	if extension == "" {
		// Если расширение пустое, добавляем его вручную
		extension = ".txt" // или другое расширение по умолчанию
	}
	fileName += extension
	logger.Println("Имя файла сгенерировано успешно")

	fileOut, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "Ошибка создания файла", http.StatusInternalServerError)
		logger.Println("Ошибка создания файла")
		return
	}
	logger.Println("файл для результата создан успешно")

	_, err = io.WriteString(fileOut, convertedString)
	defer fileOut.Close()
	if err != nil {
		http.Error(w, "Ошибка записи в файл", http.StatusInternalServerError)
		logger.Println("Ошибка записи в файл")
		return
	}
	logger.Println("Запись результата прошла успешно")

	// 7. Возврат результата
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Файл %s успешно конвертирован и сохранен как %s\n", fileHeader.Filename, fileName)
	fmt.Fprintf(w, "Результат: %s", convertedString)
}
