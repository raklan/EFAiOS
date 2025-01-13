package Engine

import "strings"

//Strips the "map_" prefix and ".json" suffix off of a map id
func StripMapId(input string) string {
	// Remove "map_" prefix
	if strings.HasPrefix(input, "map_") {
		input = input[4:]
	}
	// Remove ".json" suffix
	if strings.HasSuffix(input, ".json") {
		input = input[:len(input)-5]
	}
	return input
}
