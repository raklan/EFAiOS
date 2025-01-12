package Models

type GameMap struct {
	Id     string           `json:"id"`
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
	Row  int `json:"Row"`
	Col  int `json:"Col"`
	Type int `json:"type"`
}
