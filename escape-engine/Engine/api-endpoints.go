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
	case http.MethodGet:
		getMap(w, r)
	default:
		http.Error(w, "Method not allowed or implemented", http.StatusMethodNotAllowed)
	}
}

func AllMaps(w http.ResponseWriter, r *http.Request) {
	//getNames := r.URL.Query().Get("getNames")

	mapIds, err := GetAllMaps()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapIds)
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

func getMap(w http.ResponseWriter, r *http.Request) {
	mapId := r.URL.Query().Get("id")
	if mapId == "" {
		http.Error(w, "No map id provided", http.StatusBadRequest)
	}

	requestedMap, err := GetMapFromDB(mapId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestedMap)
}
