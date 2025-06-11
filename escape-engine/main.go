package main

import (
	"escape-engine/Models"
	"fmt"
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
	// gameState := Models.GameState{
	// 	Players: []Models.Player{
	// 		{
	// 			Id:   "1",
	// 			Name: "Ryan",
	// 			Team: Models.PlayerTeam_Human,
	// 		},
	// 		{
	// 			Id:   "4",
	// 			Name: "Meghan",
	// 			Team: Models.PlayerTeam_Human,
	// 		},
	// 		{
	// 			Id:   "2",
	// 			Name: "Christina",
	// 			Team: Models.PlayerTeam_Alien,
	// 		},
	// 		{
	// 			Id:   "3",
	// 			Name: "Blake",
	// 			Team: Models.PlayerTeam_Alien,
	// 		},
	// 	},
	// 	GameConfig: Models.GameConfig{
	// 		NumHumans: 2,
	// 		NumAliens: 2,
	// 	},
	// }

	// activeRoles := map[string]int{
	// 	Models.Role_Captain:     2,
	// 	Models.Role_Pilot:       2,
	// 	Models.Role_Soldier:     2,
	// 	Models.Role_SpeedyAlien: 2,
	// 	Models.Role_BlinkAlien:  2,
	// 	Models.Role_SilentAlien: 2,
	// }

	// requiredRoles := map[string]int{}

	// Engine.AssignRoles(&gameState, activeRoles, requiredRoles)

	// fmt.Println(gameState)

	se := Models.NewAdrenalineSurge()

	se.AddUse()

	fmt.Println(se)
}
