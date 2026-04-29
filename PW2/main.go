package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// Структура для показників системи
type Indicator struct {
	Name  string `json:"name"`  // назва показника
	Value string `json:"value"` // значення показника
	Unit  string `json:"unit"`  // одиниця вимірювання
}

// Структура для JSON-даних
type Data struct {
	Object     string      `json:"object"`     // назва об'єкта
	Indicators []Indicator `json:"indicators"` // масив показників
}

// Структура для передачі даних у шаблон
type PageData struct {
	Title      string      // заголовок сторінки
	Timestamp  string      // час перевірки
	ObjectType string      // тип об'єкта
	Indicators []Indicator // список показників
}

// Обробник HTTP-запитів
func statusHandler(w http.ResponseWriter, r *http.Request) {
	// Логування кожного запиту: метод, шлях та IP клієнта
	log.Printf("Запит: %s %s від %s", r.Method, r.URL.Path, r.RemoteAddr)

	// Вибір JSON-файлу залежно від параметра ?object=
	fileName := "data.json" // дефолтний варіант (Енергоблок)
	switch r.URL.Query().Get("object") {
	case "energy":
		fileName = "data_energy.json" // якщо ?object=energy - Система енергоспоживання
	case "dc":
		fileName = "data_dc.json" // якщо ?object=dc - Панель керування ЦОТ
	}

	// Відкриття JSON-файлу
	file, err := os.Open(fileName)
	if err != nil {
		// Якщо файл не відкрився — повертаємо помилку клієнту
		http.Error(w, "Не вдалося відкрити JSON-файл", http.StatusInternalServerError)
		return
	}
	defer file.Close() // закриваємо файл після завершення роботи

	// Декодування JSON у структуру Data
	var jsonData Data
	if err := json.NewDecoder(file).Decode(&jsonData); err != nil {
		// Якщо помилка при декодуванні — повідомляємо клієнту
		http.Error(w, "Помилка декодування JSON", http.StatusInternalServerError)
		return
	}

	// Формування даних для шаблону
	data := PageData{
		Title:      "Система аварійного електропостачання",   // заголовок сторінки
		Timestamp:  time.Now().Format("02-01-2006 15:04:05"), // поточний час у форматі ДД-ММ-РРРР ГГ:ХХ:СС
		ObjectType: jsonData.Object,                          // назва об'єкта з JSON
		Indicators: jsonData.Indicators,                      // список показників з JSON
	}

	// Парсинг HTML-шаблону (status.html)
	tmpl := template.Must(template.ParseFiles("status.html"))

	// Виконання шаблону з передачею даних
	err = tmpl.Execute(w, data)
	if err != nil {
		// Якщо помилка при генерації сторінки — повідомляємо клієнту
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Реєстрація обробника для кореневого шляху "/"
	http.HandleFunc("/", statusHandler)

	// Повідомлення у консоль про запуск сервера
	log.Println("Сервер запущено на http://localhost:8080")

	// Запуск веб-сервера на порту 8080
	// log.Fatal завершує програму, якщо сервер не зміг стартувати
	log.Fatal(http.ListenAndServe(":8080", nil))
}
