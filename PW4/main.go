package main

import (
	"log"
	"net/http"
)

func main() {
	// CRUD для систем
	http.HandleFunc("/systems", getSystems)          // GET: список систем
	http.HandleFunc("/systems/create", createSystem) // POST: створити систему
	http.HandleFunc("/systems/update", updateSystem) // PUT: оновити систему
	http.HandleFunc("/systems/delete", deleteSystem) // DELETE: видалити систему

	// Swagger UI
	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("."))))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
