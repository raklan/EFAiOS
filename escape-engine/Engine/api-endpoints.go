package Engine

import (
	"encoding/json"
	"escape-api/LogUtil"
	"escape-engine/Models"
	"log"
	"net/http"
	"strings"
)

// Given the path, gets any data that should be rendered with that requested template, if any. Only returns an error if one occurs (i.e. no data being found is not considered an error)
func GetApiData(path string) (any, error) {
	funcLogPrefix := "==GetApiData=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Getting Api Data for path {%s}", funcLogPrefix, path)

	if strings.ToLower(path) == "/maps" {
		mapIds, err := GetAllMaps()
		return mapIds, err
	}
	return nil, nil
}

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

func RoleDescription(w http.ResponseWriter, r *http.Request) {
	roleName := r.URL.Query().Get("name")
	if roleName == "" {
		http.Error(w, "No role provided", http.StatusBadRequest)
	}

	response := struct {
		RoleName        string `json:"roleName"`
		RoleDescription string `json:"roleDescription"`
	}{
		RoleName:        roleName,
		RoleDescription: Models.RoleDescriptions[roleName],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
