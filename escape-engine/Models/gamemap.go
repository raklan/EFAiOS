package Models

type GameMap struct {
	Id     string           `json:"id"`
	Name   string           `json:"name"`
	Rows   int              `json:"rows"`
	Cols   int              `json:"cols"`
	Spaces map[string]Space `json:"spaces"`
}

const (
	Space_Wall = iota
	Space_Safe
	Space_Dangerous
	Space_Pod
	Space_HumanStart
	Space_AlienStart
)

type Space struct {
	Row  int `json:"row"`
	Col  int `json:"col"`
	Type int `json:"type"`
}
