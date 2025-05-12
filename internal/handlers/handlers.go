package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	file, _, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 4. Чтение данных из файла
	data, err := io.ReadAll(file)
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

	/*// 6. Генерация имени файла
	timestamp := time.Now().UTC().String()
	extension := filepath.Ext(handler.Filename)
	newFilename := fmt.Sprintf("converted_%s%s", timestamp, extension)
	*/

	// 7. Создание и запись в локальный файл
	//err = writeToFile(newFilename, convertedString)

	fileRes, err := os.OpenFile("fileRes", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		http.Error(w, "Ошибка создания файла", http.StatusInternalServerError)
		return
	}
	err = writeToFile("fileRes", convertedString)
	defer fileRes.Close()
	if err != nil {
		http.Error(w, "Ошибка записи в файл", http.StatusInternalServerError)
		return
	}

	// 8. Возврат результата
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Файл успешно конвертирован и сохранен как %s\n", "fileRes")
	fmt.Fprintf(w, "Результат: %s", convertedString)

}

// Функция для записи в файл
func writeToFile(filename string, content string) error {
	// Создание директории, если она не существует
	dir := filepath.Dir(filename)
	if dir != "go1fl-sprint6-tpl/Data_fail" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("ошибка создания директории: %w", err)
		}
	}

	// Запись в файл
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}

	return nil
}
