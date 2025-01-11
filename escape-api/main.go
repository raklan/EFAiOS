package main

import (
	"net/http"
)

func main() {
	registerPathHandlers()
	fs := http.FileServer(http.Dir("./escape-api/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", serveHtml)

	http.ListenAndServe(":80", nil)
}

func registerPathHandlers() {
	//http.HandleFunc("/", escape_html.Test)
}

func serveHtml(w http.ResponseWriter, r *http.Request) {
	//lp := filepath.Join("escape-api", "assets")
}
