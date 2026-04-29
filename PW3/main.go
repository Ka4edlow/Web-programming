package main

import (
    "html/template"
    "net/http"
    "strconv"
    "time"
)

var tmpl = template.Must(template.ParseFiles("form.html", "result.html", "error.html")) 
// Змінна tmpl зберігає всі HTML-шаблони, які будуть використовуватись у програмі

type Equipment struct { // Структура для збереження даних обладнання
    Name    string  // Назва обладнання
    Date    string  // Дата перевірки
    Voltage float64 // Показник напруги
    Status  string  // Статус (норма/аварія)
}

func formHandler(w http.ResponseWriter, r *http.Request) { // Функція-обробник HTTP-запитів
    switch r.Method { // Перевірка методу запиту (GET або POST)
    case http.MethodGet: // Якщо запит GET
        tmpl.ExecuteTemplate(w, "form.html", nil) // Відображаємо форму для введення даних

    case http.MethodPost: // Якщо запит POST
        if err := r.ParseForm(); err != nil { // Парсимо дані з форми
            http.Error(w, "Помилка обробки форми", http.StatusBadRequest) // Якщо помилка — повертаємо повідомлення
            return // Завершуємо виконання функції
        }

        name := r.FormValue("name")       // Отримуємо значення поля "name"
        date := r.FormValue("date")       // Отримуємо значення поля "date"
        voltageStr := r.FormValue("voltage") // Отримуємо значення поля "voltage" як рядок
        status := r.FormValue("status")   // Отримуємо значення поля "status"

        voltage, err := strconv.ParseFloat(voltageStr, 64) // Конвертуємо напругу у число з плаваючою точкою
        if err != nil || voltage <= 0 { // Якщо помилка або напруга некоректна
            tmpl.ExecuteTemplate(w, "error.html", "Некоректне значення напруги") // Виводимо повідомлення про помилку
            return // Завершуємо виконання функції
        }

        if _, err := time.Parse("2006-01-02", date); err != nil { // Перевіряємо коректність дати
            tmpl.ExecuteTemplate(w, "error.html", "Некоректна дата") // Якщо дата некоректна — повідомлення про помилку
            return // Завершуємо виконання функції
        }

        eq := Equipment{Name: name, Date: date, Voltage: voltage, Status: status} // Створюємо об'єкт Equipment з даними
        tmpl.ExecuteTemplate(w, "result.html", eq) // Відображаємо шаблон результатів з переданими даними

    default: // Якщо метод запиту не підтримується
        http.Error(w, "Метод не підтримується", http.StatusMethodNotAllowed) // Повертаємо повідомлення про помилку
    }
}

func main() { // Головна функція програми
    http.HandleFunc("/", formHandler) // Реєструємо обробник для кореневого маршруту "/"
    println("Сервер запущено на http://localhost:8080") // Виводимо повідомлення про запуск сервера у консоль
    if err := http.ListenAndServe(":8080", nil); err != nil { // Запускаємо сервер на порту 8080
        panic(err) // Якщо виникла помилка — аварійно завершуємо програму
    }
}
