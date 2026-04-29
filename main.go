package main // оголошення головного пакету програми

import (
    "fmt"      // пакет для форматованого виводу
    "log"      // пакет для логування повідомлень
    "net/http" // пакет для створення веб-сервера та роботи з HTTP
    "time"     // пакет для роботи з датою та часом
)

// statusHandler — функція-обробник HTTP-запитів
func statusHandler(w http.ResponseWriter, r *http.Request) {
    // Логування кожного запиту: метод, шлях та IP клієнта
    log.Printf("Запит: %s %s від %s", r.Method, r.URL.Path, r.RemoteAddr)

    // Отримання параметра ?mode=... (режим роботи системи)
    mode := r.URL.Query().Get("mode")
    if mode == "" {
        mode = "normal" // якщо параметр не задано — режим "normal"
    }

    // Отримання параметра ?interface=... (варіант інтерфейсу)
    ui := r.URL.Query().Get("interface")
    if ui == "" {
        ui = "default" // якщо параметр не задано — інтерфейс "default"
    }

    // Формування HTML-відповіді
    w.Header().Set("Content-Type", "text/html; charset=utf-8") // заголовок відповіді
    fmt.Fprintf(w, "<html><head><title>Система аварійного електропостачання</title></head><body>")
    fmt.Fprintf(w, "<h1>Стан системи аварійного електропостачання</h1>")
    fmt.Fprintf(w, "<p>Час перевірки: %s</p>", time.Now().Format("15:04:05")) // поточний час

    // Відображення режиму роботи системи залежно від параметра mode
    switch mode {
    case "emergency":
        fmt.Fprintf(w, "<p style='color:red;'> Аварія! Перехід на резервне живлення.</p>")
    case "backup":
        fmt.Fprintf(w, "<p style='color:orange;'> Система працює від резервного джерела.</p>")
    default:
        fmt.Fprintf(w, "<p style='color:green;'> Система працює у штатному режимі.</p>")
    }

    // Варіативність інтерфейсу залежно від параметра interface
    fmt.Fprintf(w, "<hr>") // горизонтальна лінія для відділення блоків
    switch ui {
    case "monitoring":
        fmt.Fprintf(w, "<h2>Інтерфейс: Станція моніторингу</h2>")
        fmt.Fprintf(w, "<p>Відображення датчиків напруги та частоти.</p>")
    case "energy":
        fmt.Fprintf(w, "<h2>Інтерфейс: Система енергоспоживання</h2>")
        fmt.Fprintf(w, "<p>Графік споживання електроенергії та баланс резервів.</p>")
    case "dc":
        fmt.Fprintf(w, "<h2>Інтерфейс: Панель керування ЦОТ</h2>")
        fmt.Fprintf(w, "<p>Управління живленням серверів та UPS.</p>")
    default:
        fmt.Fprintf(w, "<h2>Інтерфейс: Станція моніторингу (за замовчуванням)</h2>")
    }
    fmt.Fprintf(w, "</body></html>") // завершення HTML-документа
}

func main() {
    // Реєстрація обробника для кореневого шляху "/"
    http.HandleFunc("/", statusHandler)

    // Повідомлення у консоль про запуск сервера
    fmt.Println("Сервер запущено на http://localhost:8080")

    // Запуск веб-сервера на порту 8080, з обробкою помилок через log.Fatal
    log.Fatal(http.ListenAndServe(":8080", nil))
}
