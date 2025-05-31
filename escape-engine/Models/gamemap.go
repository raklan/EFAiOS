package Models

import (
	"fmt"
	"maps"
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

func (gameMap GameMap) GetSpacesWithinNthAdjacency(n int, homeSpaceKey string) map[string]Space {
	spaces := map[string]Space{}

	homeSpace := gameMap.Spaces[homeSpaceKey]

	rowNum, colNum := homeSpace.GetRowAsInt(), homeSpace.Col

	//Get each of the 6 directions. Need to handle even/odd columns differently because of how hex maps line up

	if colNum%2 == 0 { //Even column
		//Up-Left
		spaceKey := fmt.Sprintf("%s-%d", GetRowAsLetter(rowNum-1), colNum-1)
		if space, exists := gameMap.Spaces[spaceKey]; exists {
			spaces[spaceKey] = space
		}

		//Up
		spaceKey = fmt.Sprintf("%s-%d", GetRowAsLetter(rowNum-1), colNum)
		if space, exists := gameMap.Spaces[spaceKey]; exists {
			spaces[spaceKey] = space
		}

		//Up-Right
		spaceKey = fmt.Sprintf("%s-%d", GetRowAsLetter(rowNum-1), colNum+1)
		if space, exists := gameMap.Spaces[spaceKey]; exists {
			spaces[spaceKey] = space
		}

		//Down-Left
		spaceKey = fmt.Sprintf("%s-%d", GetRowAsLetter(rowNum), colNum-1)
		if space, exists := gameMap.Spaces[spaceKey]; exists {
			spaces[spaceKey] = space
		}

		//Down
		spaceKey = fmt.Sprintf("%s-%d", GetRowAsLetter(rowNum+1), colNum)
		if space, exists := gameMap.Spaces[spaceKey]; exists {
			spaces[spaceKey] = space
		}

		//Down-Right
		spaceKey = fmt.Sprintf("%s-%d", GetRowAsLetter(rowNum), colNum+1)
		if space, exists := gameMap.Spaces[spaceKey]; exists {
			spaces[spaceKey] = space
		}
	} else {
		//Up-Left
		spaceKey := fmt.Sprintf("%s-%d", GetRowAsLetter(rowNum), colNum-1)
		if space, exists := gameMap.Spaces[spaceKey]; exists {
			spaces[spaceKey] = space
		}

		//Up
		spaceKey = fmt.Sprintf("%s-%d", GetRowAsLetter(rowNum-1), colNum)
		if space, exists := gameMap.Spaces[spaceKey]; exists {
			spaces[spaceKey] = space
		}

		//Up-Right
		spaceKey = fmt.Sprintf("%s-%d", GetRowAsLetter(rowNum), colNum+1)
		if space, exists := gameMap.Spaces[spaceKey]; exists {
			spaces[spaceKey] = space
		}

		//Down-Left
		spaceKey = fmt.Sprintf("%s-%d", GetRowAsLetter(rowNum+1), colNum-1)
		if space, exists := gameMap.Spaces[spaceKey]; exists {
			spaces[spaceKey] = space
		}

		//Down
		spaceKey = fmt.Sprintf("%s-%d", GetRowAsLetter(rowNum+1), colNum)
		if space, exists := gameMap.Spaces[spaceKey]; exists {
			spaces[spaceKey] = space
		}

		//Down-Right
		spaceKey = fmt.Sprintf("%s-%d", GetRowAsLetter(rowNum+1), colNum+1)
		if space, exists := gameMap.Spaces[spaceKey]; exists {
			spaces[spaceKey] = space
		}
	}

	neighbors := maps.Clone(spaces) //Make a clone to iterate over the neighbors to avoid the collection changing while iterating over it

	if n > 1 {
		for neighborKey := range neighbors {
			maps.Copy(spaces, gameMap.GetSpacesWithinNthAdjacency(n-1, neighborKey))
			delete(spaces, homeSpaceKey)
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

func GetMapKey(row string, col int) string {
	return fmt.Sprintf("%s-%d", row, col)
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
