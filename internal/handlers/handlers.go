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
	// 1. Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		log.Println("Метод не поддерживается")
		return
	}

	// 2. Парсинг формы
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB ограничение
		http.Error(w, "Ошибка парсинга формы", http.StatusInternalServerError)
		return
	}

	// 3. Получение файла из формы
	fileHtml, fileHeader, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusInternalServerError)
		return
	}
	defer fileHtml.Close()

	// 4. Чтение данных из файла
	data, err := io.ReadAll(fileHtml)
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	// 5. Автоопределение и конвертация
	convertedString, err := service.Service(string(data))
	if err != nil {
		http.Error(w, "Ошибка конвертации", http.StatusInternalServerError)
		return
	}

	// 6. Генерация имени файла, создание и запись в локальный файл
	fileName := time.Now().UTC().Format("2006-01-02_15-04-05")
	extension := filepath.Ext(fileHeader.Filename)
	// Проверяем, что расширение не пустое
	if extension == "" {
		// Если расширение пустое, добавляем его вручную
		extension = ".txt" // или другое расширение по умолчанию
	}
	fileName += extension

	fileOut, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "Ошибка создания файла", http.StatusInternalServerError)
		return
	}
	_, err = io.WriteString(fileOut, convertedString)
	defer fileOut.Close()
	if err != nil {
		http.Error(w, "Ошибка записи в файл", http.StatusInternalServerError)
		return
	}

	// 7. Возврат результата
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Файл %s успешно конвертирован и сохранен как %s\n", fileHeader.Filename, fileName)
	fmt.Fprintf(w, "Результат: %s", convertedString)

}
