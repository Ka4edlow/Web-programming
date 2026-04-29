package main

import (
    "encoding/json"
    "net/http"
    "strconv"
)

// Структура системи
type PowerSystem struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    Source string `json:"source"` // зелена електрика / теплова (ТЕЦ)
    State  string `json:"state"`  // стабільний / критичний / аварійний
}

// Масив для збереження систем
var systems = []PowerSystem{}
var nextID = 1

// GET /systems – отримати всі системи
func getSystems(w http.ResponseWriter, r *http.Request) {
    respondWithJSON(w, http.StatusOK, systems)
}

// POST /systems/create – створити нову систему
func createSystem(w http.ResponseWriter, r *http.Request) {
    var sys PowerSystem
    if err := json.NewDecoder(r.Body).Decode(&sys); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }
    sys.ID = nextID
    nextID++
    systems = append(systems, sys)
    respondWithJSON(w, http.StatusCreated, sys)
}

// PUT /systems/update?id=1 – оновити систему
func updateSystem(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid ID")
        return
    }
    var updated PowerSystem
    if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }
    for i, s := range systems {
        if s.ID == id {
            updated.ID = id
            systems[i] = updated
            respondWithJSON(w, http.StatusOK, updated)
            return
        }
    }
    respondWithError(w, http.StatusNotFound, "System not found")
}

// DELETE /systems/delete?id=1 – видалити систему
func deleteSystem(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid ID")
        return
    }
    for i, s := range systems {
        if s.ID == id {
            systems = append(systems[:i], systems[i+1:]...)
            respondWithJSON(w, http.StatusOK, map[string]string{"message": "Deleted"})
            return
        }
    }
    respondWithError(w, http.StatusNotFound, "System not found")
}

// Допоміжні функції
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}
