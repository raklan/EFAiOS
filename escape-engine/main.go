package main

import (
	"os"
	"time"
)

// func main() {
// 	gameMap := Models.GameMap{
// 		Rows: 3,
// 		Cols: 5,
// 		Spaces: map[string]Models.Space{
// 			"A-0": Models.Space{
// 				Row: "A",
// 				Col: 0,
// 			},
// 			"A-1": Models.Space{
// 				Row: "A",
// 				Col: 1,
// 			},
// 			"A-2": Models.Space{
// 				Row: "A",
// 				Col: 2,
// 			},
// 			"A-3": Models.Space{
// 				Row: "A",
// 				Col: 3,
// 			},
// 			"A-4": Models.Space{
// 				Row: "A",
// 				Col: 4,
// 			},
// 			"B-0": Models.Space{
// 				Row: "B",
// 				Col: 0,
// 			},
// 			"B-1": Models.Space{
// 				Row: "B",
// 				Col: 1,
// 			},
// 			"B-2": Models.Space{
// 				Row: "B",
// 				Col: 2,
// 			},
// 			"B-3": Models.Space{
// 				Row: "B",
// 				Col: 3,
// 			},
// 			"B-4": Models.Space{
// 				Row: "B",
// 				Col: 4,
// 			},
// 			"C-0": Models.Space{
// 				Row: "C",
// 				Col: 0,
// 			},
// 			"C-1": Models.Space{
// 				Row: "C",
// 				Col: 1,
// 			},
// 			"C-2": Models.Space{
// 				Row: "C",
// 				Col: 2,
// 			},
// 			"C-3": Models.Space{
// 				Row: "C",
// 				Col: 3,
// 			},
// 			"C-4": Models.Space{
// 				Row: "C",
// 				Col: 4,
// 			},
// 			"D-0": Models.Space{
// 				Row: "C",
// 				Col: 0,
// 			},
// 			"D-1": Models.Space{
// 				Row: "C",
// 				Col: 1,
// 			},
// 			"D-2": Models.Space{
// 				Row: "C",
// 				Col: 2,
// 			},
// 			"D-3": Models.Space{
// 				Row: "C",
// 				Col: 3,
// 			},
// 			"D-4": Models.Space{
// 				Row: "C",
// 				Col: 4,
// 			},
// 			"E-0": Models.Space{
// 				Row: "C",
// 				Col: 0,
// 			},
// 			"E-1": Models.Space{
// 				Row: "C",
// 				Col: 1,
// 			},
// 			"E-2": Models.Space{
// 				Row: "C",
// 				Col: 2,
// 			},
// 			"E-3": Models.Space{
// 				Row: "C",
// 				Col: 3,
// 			},
// 			"E-4": Models.Space{
// 				Row: "C",
// 				Col: 4,
// 			},
// 		},
// 	}

// 	res := Actions.GetSpacesWithinNthAdjacency(3, "A-1", gameMap)
// 	fmt.Println(res)
// }

func main() {
	files, _ := os.ReadDir("./maps")
	for _, file := range files {
		stats, _ := os.Stat("./maps/" + file.Name())
		expirationTime := stats.ModTime().AddDate(0, 1, 0)
		if time.Now().After(expirationTime) {
			os.Remove("./maps/" + file.Name())
		}
	}
	// stats, _ := os.Stat("./maps/map_1738255395028DWVM9PW8T1.json")
	// fmt.Println(time.Now())
	// expirationTime := stats.ModTime().AddDate(0, 1, 0)
	// fmt.Println(stats.ModTime(), expirationTime)
	// fmt.Println(time.Now().Before(expirationTime))
}
