package Models

import (
	"fmt"
	"math"
)

type GameMap struct {
	Id     string           `json:"id"`
	Name   string           `json:"name"`
	Rows   int              `json:"rows"`
	Cols   int              `json:"cols"`
	Spaces map[string]Space `json:"spaces"`
}

func (gameMap *GameMap) GetSpacesOfType(spaceType int) []Space {
	spaces := []Space{}
	for _, space := range gameMap.Spaces {
		if space.Type == spaceType {
			spaces = append(spaces, space)
		}
	}
	return spaces
}

const (
	Space_Wall = iota
	Space_Safe
	Space_Dangerous
	Space_Pod
	Space_UsedPod
	Space_HumanStart
	Space_AlienStart
)

type Space struct {
	Row  string `json:"row"`
	Col  int    `json:"col"`
	Type int    `json:"type"`
}

func (space Space) GetMapKey() string {
	return fmt.Sprintf("%s-%d", space.Row, space.Col)
}

// Gets the space's row as an int. A = 0, B = 1, ... Z = 25, AA = 26, etc
func (space Space) GetRowAsInt() int {
	row := 0

	for i := range space.Row {
		row += int(space.Row[i]) - 65 + (i * 26)
	}

	return row
}

func GetRowAsLetter(rowNum int) string {
	letterCode := ""

	for rowNum >= 0 {
		letterCode += string(rune(65 + rowNum%26))
		rowNum = int(math.Floor(float64(rowNum)/26)) - 1
	}

	return letterCode
}
