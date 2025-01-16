package Engine

import "strings"

//Strips the "map_" prefix and ".json" suffix off of a map id
func StripMapId(input string) string {
	// Remove "map_" prefix
	input = strings.TrimPrefix(input, "map_")
	// Remove ".json" suffix
	input = strings.TrimSuffix(input, ".json")
	return input
}
