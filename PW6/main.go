package main

import (
	"log"
	"net/http"
)

func main() {
	initDB()

	// API маршрути
	http.HandleFunc("/systems", getSystems)          // GET
	http.HandleFunc("/systems/create", createSystem) // POST
	http.HandleFunc("/systems/update", updateSystem) // PUT
	http.HandleFunc("/systems/delete", deleteSystem) // DELETE

	// Маршрути для користувачів
	http.HandleFunc("/register", registerUser) // POST
	http.HandleFunc("/login", loginUser)       // POST
	http.HandleFunc("/profile", getProfile)    // GET

	// Swagger UI (документація API)
	http.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger"))))

	// Статичний інтерфейс
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("./static"))))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
