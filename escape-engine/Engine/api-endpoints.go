package Engine

import (
	"encoding/json"
	"escape-engine/Models"
	"net/http"
)

func Map(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		saveMap(w, r)
	default:
		http.Error(w, "Method not allowed or implemented", http.StatusMethodNotAllowed)
	}
}

func saveMap(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	req := Models.GameMap{}

	err := d.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	saved, err := SaveMapToDB(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Return the game data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saved)
}
