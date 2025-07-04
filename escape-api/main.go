package main

import (
	"escape-engine/Engine"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	setUpLogging()
	Engine.PrepareFilesystem()
	startServer()
}

func startServer() {
	log.Println("=========================Starting Server========================")
	fmt.Println("=========================Starting Server========================")
	fs := http.FileServer(http.Dir("./escape-api/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.Handle("/favicon.ico", fs)

	http.HandleFunc("/", serveHtml)

	http.HandleFunc("/api/map", Engine.Map)
	http.HandleFunc("/api/allMaps", Engine.AllMaps)
	http.HandleFunc("/api/mapForLobby", Engine.GetMapForLobby)

	http.HandleFunc("/api/role", Engine.RoleDescription)

	http.HandleFunc("/lobby/host", Engine.HostLobby)
	http.HandleFunc("/lobby/join", Engine.HandleJoinLobby)
	http.HandleFunc("/lobby/rejoin", Engine.HandleRejoinLobby)

	fmt.Printf("Server has started listening on port 80. Connect to %s from a web browser to play!\n", GetLocalIP())
	http.ListenAndServe(":80", nil)
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
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

	templateData, err := Engine.GetApiData(r.URL.Path, r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temp := template.New("layout.html").Funcs(template.FuncMap{
		"StripMapId":         Engine.StripMapId,
		"GetMapName":         Engine.GetMapName,
		"GetMapDescription":  Engine.GetMapDescription,
		"GetConfigPresets":   Engine.GetConfigPresets,
		"NumPlayersForRecap": Engine.NumPlayers,
		"ToLowercase":        strings.ToLower,
		"EqualZero":          Engine.EqualZero,
	})

	var tmpl *template.Template
	if strings.Contains(strings.ToLower(r.URL.Path), "compendium") {
		compendiumPath := filepath.Join("escape-api", "assets", "html", "templates", "layout_compendium.html")
		tmpl, err = temp.ParseFiles(layoutPath, compendiumPath, requestedFilePath)
	} else {
		tmpl, err = temp.ParseFiles(layoutPath, requestedFilePath)
	}
	if err != nil {
		tmpl, err = template.ParseFiles(layoutPath, filepath.Join("escape-api", "assets", "html", "index.html"))
		if err != nil {
			fmt.Fprintf(w, "It broke")
			return
		}
	}

	tmpl.ExecuteTemplate(w, "layout", templateData)
}
