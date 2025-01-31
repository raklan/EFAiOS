package Models

import (
	"fmt"
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
