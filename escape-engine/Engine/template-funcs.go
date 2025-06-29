package Engine

import (
	"escape-engine/Models/GameConfig"
	"escape-engine/Models/Recap"
	"html/template"
	"strings"
)

// Strips the "map_" prefix and ".json" suffix off of a map id
func StripMapId(input string) string {
	// Remove "map_" prefix
	input = strings.TrimPrefix(input, "map_")
	// Remove ".json" suffix
	input = strings.TrimSuffix(input, ".json")
	return input
}

// Extracts the Name of the map from a map key
func GetMapName(input string) string {
	gameMap, err := GetMapFromDB(StripMapId(input))
	if err != nil {
		return "Error finding Map"
	}

	return gameMap.Name
}

func GetMapDescription(input string) template.HTML { //TODO: Sanitize description somehow to prevent XSS?
	gameMap, err := GetMapFromDB(StripMapId(input))
	if err != nil {
		return "Error finding Map"
	}

	return template.HTML(gameMap.Description)
}

func GetConfigPresets() []GameConfig.GameConfigPreset {
	return GameConfig.GetConfigPresets()
}

func NumPlayers(players []Recap.PlayerRecap) int {
	return len(players)
}
