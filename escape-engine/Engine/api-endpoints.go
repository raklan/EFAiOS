package Engine

import (
	"encoding/json"
	"escape-api/LogUtil"
	"escape-engine/Models"
	"escape-engine/Models/Recap"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Given the path, gets any data that should be rendered with that requested template, if any. Only returns an error if one occurs (i.e. no data being found is not considered an error)
func GetApiData(path string, query url.Values) (any, error) {
	funcLogPrefix := "==GetApiData=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Getting Api Data for path {%s}", funcLogPrefix, path)

	if strings.ToLower(path) == "/maps" {
		mapIds, err := GetAllMaps()
		return mapIds, err
	} else if strings.ToLower(path) == "/recap" {
		recap := GetRecap(query)
		return recap, nil
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

func GetRecap(query url.Values) Recap.Recap {
	roomCode := query.Get("roomCode")

	lobby, err := GetLobbyFromFs(roomCode)
	if err != nil {
		return Recap.Recap{}
	}

	if lobby.Status != Models.LobbyStatus_Ended {
		return Recap.Recap{
			MapName: "Game has not ended yet",
		}
	}

	recap, err := Recap.GetRecapFromFs(lobby.GameStateId)
	if err != nil {
		LogError("==GetRecap (API Data)==", err)
		return Recap.Recap{}
	}
	return recap
}
