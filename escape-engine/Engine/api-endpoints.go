package Engine

import "net/http"

func Map(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		saveMap(w, r)
	default:
		http.Error(w, "Method not allowed or implemented", http.StatusMethodNotAllowed)
	}
}

func saveMap(w http.ResponseWriter, r *http.Request) {

}
