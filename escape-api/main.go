package main

import (
	"escape-engine/Engine"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	startServer()
}

func startServer() {

	fs := http.FileServer(http.Dir("./escape-api/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.Handle("/favicon.ico", fs)

	http.HandleFunc("/", serveHtml)

	http.HandleFunc("/api/map", Engine.Map)
	http.HandleFunc("/api/allMaps", Engine.AllMaps)

	log.Println("=========================Starting Server========================")

	http.ListenAndServe(":80", nil)
}

func setUpLogging() {
	logName := "./logs/server.log"
	log.SetPrefix("ESCAPE-API: ")
	log.SetOutput(&lumberjack.Logger{
		Filename: logName,
		MaxSize:  1,
		MaxAge:   7,
		Compress: false,
	})
}

func serveHtml(w http.ResponseWriter, r *http.Request) {
	layoutPath := filepath.Join("escape-api", "assets", "html", "templates", "layout.html")
	requestedFilePath := filepath.Join("escape-api", "assets", "html", fmt.Sprintf("%s.html", filepath.Clean(r.URL.Path)))

	tmpl, err := template.ParseFiles(layoutPath, requestedFilePath)
	if err != nil {
		tmpl, err = template.ParseFiles(layoutPath, filepath.Join("escape-api", "assets", "html", "index.html"))
		if err != nil {
			fmt.Fprintf(w, "It broke")
			return
		}
	}
	tmpl.ExecuteTemplate(w, "layout", nil)
}
