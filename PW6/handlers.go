package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    _ "github.com/go-sql-driver/mysql"
    "golang.org/x/crypto/bcrypt"
)

// Структура PowerSystem описує модель даних для системи електропостачання.
// Вона використовується для зберігання даних з БД та передачі у JSON.
type PowerSystem struct {
	ID     int    `json:"id"`     // Унікальний ідентифікатор системи
	Name   string `json:"name"`   // Назва системи
	Source string `json:"source"` // Джерело електропостачання
	State  string `json:"state"`  // Стан системи
}

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
}

// Глобальна змінна для з'єднання з базою даних
var db *sql.DB

// Функція initDB встановлює з'єднання з MySQL.
// Використовується DSN (Data Source Name) із параметрами кодування та локалізації.
func initDB() {
	var err error
	dsn := "poweruser:powerpass@tcp(127.0.0.1:3306)/powerdb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = sql.Open("mysql", dsn) // Відкриття з'єднання
	if err != nil {
		panic(err) // Якщо помилка — аварійне завершення
	}
	if err = db.Ping(); err != nil {
		panic(err) // Перевірка доступності БД
	}
	log.Println("Connected to MySQL") // Лог повідомлення про успішне підключення
}

// Метод getSystems виконує SELECT-запит до таблиці systems.
// Він повертає всі записи у форматі JSON.
func getSystems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, source, state FROM systems")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close() // Закриваємо результати після використання

	var systems []PowerSystem
	for rows.Next() {
		var sys PowerSystem
		// Зчитуємо дані з рядка у структуру
		if err := rows.Scan(&sys.ID, &sys.Name, &sys.Source, &sys.State); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		systems = append(systems, sys)
	}
	respondWithJSON(w, http.StatusOK, systems) // Відправляємо JSON-відповідь
}

// Метод createSystem додає новий запис у таблицю systems.
// Він приймає JSON із клієнта, виконує INSERT і повертає створений запис.
func createSystem(w http.ResponseWriter, r *http.Request) {
	var sys PowerSystem
	// Декодуємо JSON із тіла запиту у структуру
	if err := json.NewDecoder(r.Body).Decode(&sys); err != nil {
		log.Println("JSON decode error:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	log.Println("Trying to insert:", sys.Name, sys.Source, sys.State)

	// Виконуємо SQL-операцію INSERT
	res, err := db.Exec("INSERT INTO systems (name, source, state) VALUES (?, ?, ?)",
		sys.Name, sys.Source, sys.State)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error: "+err.Error())
		return
	}
	// Отримуємо ID нового запису
	id, err := res.LastInsertId()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot get ID: "+err.Error())
		return
	}

	sys.ID = int(id)
	respondWithJSON(w, http.StatusCreated, sys) // Повертаємо створений запис
}

// Метод updateSystem оновлює існуючий запис у таблиці systems.
// Він приймає ID через параметр URL і нові дані через JSON.
func updateSystem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // Отримуємо ID із параметра URL
	id, err := strconv.Atoi(idStr)   // Конвертуємо у число
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	var sys PowerSystem
	// Декодуємо JSON із тіла запиту
	if err := json.NewDecoder(r.Body).Decode(&sys); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	// Виконуємо SQL-операцію UPDATE
	_, err = db.Exec("UPDATE systems SET name=?, source=?, state=? WHERE id=?",
		sys.Name, sys.Source, sys.State, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sys.ID = id
	respondWithJSON(w, http.StatusOK, sys) // Повертаємо оновлений запис
}

// Метод deleteSystem видаляє запис із таблиці systems за ID.
func deleteSystem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // Отримуємо ID із параметра URL
	id, err := strconv.Atoi(idStr)   // Конвертуємо у число
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	// Виконуємо SQL-операцію DELETE
	_, err = db.Exec("DELETE FROM systems WHERE id=?", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Повертаємо повідомлення про успішне видалення
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Deleted"})
}

// Допоміжна функція для формування JSON-відповіді.
// Встановлює заголовки, статус і серіалізує payload у JSON.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// Допоміжна функція для формування JSON-помилки.
// Використовує respondWithJSON для повернення повідомлення про помилку.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func registerUser(w http.ResponseWriter, r *http.Request) {
    var u User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }
    hashed, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    _, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", u.Username, string(hashed))
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "DB error: "+err.Error())
        return
    }
    respondWithJSON(w, http.StatusCreated, map[string]string{"message": "User registered"})
}

func loginUser(w http.ResponseWriter, r *http.Request) {
    var u User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }
    var hashed string
    err := db.QueryRow("SELECT password FROM users WHERE username=?", u.Username).Scan(&hashed)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
        return
    }
    if bcrypt.CompareHashAndPassword([]byte(hashed), []byte(u.Password)) != nil {
        respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
        return
    }
    // тут можна створити cookie або JWT
    http.SetCookie(w, &http.Cookie{Name: "session_user", Value: u.Username, Path: "/"})
    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Login successful"})
}

func getProfile(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_user")
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Not logged in")
        return
    }
    respondWithJSON(w, http.StatusOK, map[string]string{"username": cookie.Value})
}