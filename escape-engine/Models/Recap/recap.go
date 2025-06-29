package Recap

import (
	"encoding/json"
	"escape-api/LogUtil"
	"fmt"
	"log"
	"os"
	"slices"
	"sync"
	"time"
)

type Recap struct {
	GameStateId string        `json:"gameStateId"`
	MapName     string        `json:"mapName"`
	Players     []PlayerRecap `json:"players"`
}

type PlayerRecap struct {
	PlayerId   string         `json:"playerId"`
	PlayerName string         `json:"playerName"`
	PlayerTeam string         `json:"playerTeam"`
	PlayerRole string         `json:"playerRole"`
	Turns      map[int]string `json:"turns"`
}

var recapFileMutex = sync.Mutex{}

func SaveRecapToFs(recap Recap) (Recap, error) {
	funcLogPrefix := "==SaveRecapToFs=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix("ESCAPE-ENGINE", "Recap")
	asJson, err := json.Marshal(recap)
	if err != nil {
		return recap, err
	}

	filename := "recap_" + recap.GameStateId + ".json"

	f, err := os.Create(fmt.Sprintf("./recaps/%s", filename))
	if err != nil {
		log.Printf("%s Ran into unrecoverable error trying to save recap to filesystem", funcLogPrefix)
		f.Close()
		return recap, err
	}
	_, err = f.Write(asJson)
	f.Close()
	if err != nil {
		return recap, err
	}

	//Kick off goroutine clearing out unused recaps
	go clearOutOldFiles("./recaps/")

	return recap, nil
}

func GetRecapFromFs(gameStateId string) (Recap, error) {
	funcLogPrefix := "==GetRecapFromFs=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix("ESCAPE-ENGINE", "Recap")

	log.Printf("%s Getting recap from FS with GameStateId == {%s}", funcLogPrefix, gameStateId)
	data, err := os.ReadFile(fmt.Sprintf("./recaps/recap_%s.json", gameStateId))
	if err != nil {
		log.Printf("ERROR! %s", err)
		return Recap{}, err
	}

	parsed := Recap{}

	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return parsed, err
	}

	return parsed, nil
}

func AddDataToRecap(gameStateId string, playerId string, turnNumber int, turnDescription string) {
	recapFileMutex.Lock()
	recap, err := GetRecapFromFs(gameStateId)
	if err != nil {
		panic(fmt.Sprintf("Couldn't find recap with gameStateId: %s", gameStateId))
	}

	if indexToAddTo := slices.IndexFunc(recap.Players, func(p PlayerRecap) bool { return p.PlayerId == playerId }); indexToAddTo != -1 {
		if recap.Players[indexToAddTo].Turns[turnNumber] != "" {
			recap.Players[indexToAddTo].Turns[turnNumber] += fmt.Sprintf(", %s", turnDescription)
		} else {
			recap.Players[indexToAddTo].Turns[turnNumber] = turnDescription
		}

		SaveRecapToFs(recap)
	}
	recapFileMutex.Unlock()
}

func clearOutOldFiles(directory string) {
	files, _ := os.ReadDir(directory)
	for _, file := range files {
		fullFileName := directory + file.Name()
		stats, _ := os.Stat(fullFileName)
		expirationTime := stats.ModTime().AddDate(0, 0, 7)
		if time.Now().After(expirationTime) {
			os.Remove(fullFileName)
		}
	}
}
