package GameConfig

import "encoding/json"

type GameConfigPreset struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ConfigJson  string `json:"configJson"`
}

// GameState-specific config as defined by the Host
type GameConfig struct {
	//Number of Humans currently in the Game. The Game automatically ends when this number hits 0.
	NumHumans int `json:"numHumans"`
	//Number of Aliens currently in the Game
	NumAliens int `json:"numAliens"`
	//Number of Working Escape Pods left. The Game automatically ends when this number hits 0.
	NumWorkingPods int `json:"numWorkingPods"`
	//Number of Broken Escape Pods left
	NumBrokenPods int `json:"numBrokenPods"`
	//Number of turns before the game should end
	NumTurns int `json:"numTurns"`
	//Whether Alien players should join the Spectators team upon death
	AliensRespawn bool `json:"aliensRespawn"`
	//Whether players' turns should automatically end after moving (and choosing to attack, for aliens) if there are no cards in their hand
	AutoTurnEnd bool `json:"autoTurnEnd"`
	//Which cards should be active, as well as how many of each
	ActiveCards map[string]int `json:"activeCards"`
	//Which roles should be active, as well as the maximum number allowed to be present. Should be >= that role's presence in RequiredRoles, if it's required
	ActiveRoles map[string]int `json:"activeRoles"`
	//Which roles should be guaranteed to be in the game, as well as the number of players that should have that role
	RequiredRoles map[string]int `json:"requiredRoles"`
	//Which StatusEffects should be active, as well as their priority
	ActiveStatusEffects map[string]int `json:"activeStatusEffects"`
}

func GetConfigPresets() []GameConfigPreset {
	return []GameConfigPreset{
		{
			Name:        "Tabletop",
			Description: "The Classic experience, with all settings exactly as they appear in the Tabletop version",
			ConfigJson: GetConfigAsJsonString(GameConfig{
				NumHumans:      0,
				NumAliens:      0,
				NumWorkingPods: 4,
				NumBrokenPods:  1,
				NumTurns:       40,
				AliensRespawn:  false,
				AutoTurnEnd:    false,
				ActiveCards: map[string]int{
					"Red Card":   24,
					"Green Card": 26,
					"White Card": 4,
					"Adrenaline": 3,
					"Attack":     1,
					"Cat":        2,
					"Clone":      1,
					"Defense":    1,
					"Mutation":   1,
					"Sedatives":  1,
					"Sensor":     1,
					"Spotlight":  2,
					"Teleport":   1,
				},
				ActiveRoles: map[string]int{
					"Captain":           1,
					"Pilot":             1,
					"Copilot":           1,
					"Soldier":           1,
					"Psychologist":      1,
					"Executive Officer": 1,
					"Medic":             1,
					"Engineer":          1,
					"Fast":              1,
					"Surge":             1,
					"Blink":             1,
					"Silent":            1,
					"Brute":             1,
					"Invisible":         1,
					"Lurking":           1,
					"Psychic":           1,
				},
				RequiredRoles: map[string]int{
					"Captain":           0,
					"Pilot":             0,
					"Copilot":           0,
					"Soldier":           0,
					"Psychologist":      0,
					"Executive Officer": 0,
					"Medic":             0,
					"Engineer":          0,
					"Fast":              0,
					"Surge":             0,
					"Blink":             0,
					"Silent":            0,
					"Brute":             0,
					"Invisible":         0,
					"Lurking":           0,
					"Psychic":           0,
				},
			}),
		},
		{
			Name:        "The Basics",
			Description: "A simplified version of the Tabletop preset, with all Roles turned off and every Item card replaced with a White \"Silent\" Card. Ideal for people just learning the game",
			ConfigJson: GetConfigAsJsonString(GameConfig{
				NumHumans:      0,
				NumAliens:      0,
				NumWorkingPods: 4,
				NumBrokenPods:  1,
				NumTurns:       40,
				AliensRespawn:  false,
				AutoTurnEnd:    false,
				ActiveCards: map[string]int{
					"Red Card":   24,
					"Green Card": 26,
					"White Card": 18,
					"Adrenaline": 0,
					"Attack":     0,
					"Cat":        0,
					"Clone":      0,
					"Defense":    0,
					"Mutation":   0,
					"Sedatives":  0,
					"Sensor":     0,
					"Spotlight":  0,
					"Teleport":   0,
				},
				ActiveRoles: map[string]int{
					"Captain":           0,
					"Pilot":             0,
					"Copilot":           0,
					"Soldier":           0,
					"Psychologist":      0,
					"Executive Officer": 0,
					"Medic":             0,
					"Engineer":          0,
					"Fast":              0,
					"Surge":             0,
					"Blink":             0,
					"Silent":            0,
					"Brute":             0,
					"Invisible":         0,
					"Lurking":           0,
					"Psychic":           0,
				},
				RequiredRoles: map[string]int{
					"Captain":           0,
					"Pilot":             0,
					"Copilot":           0,
					"Soldier":           0,
					"Psychologist":      0,
					"Executive Officer": 0,
					"Medic":             0,
					"Engineer":          0,
					"Fast":              0,
					"Surge":             0,
					"Blink":             0,
					"Silent":            0,
					"Brute":             0,
					"Invisible":         0,
					"Lurking":           0,
					"Psychic":           0,
				},
			}),
		},
		{
			Name:        "Feeding Frenzy",
			Description: "Conner's brutal preset. Can you escape the Aliens when the only role is Psychologist?",
			ConfigJson: GetConfigAsJsonString(GameConfig{
				NumHumans:      0,
				NumAliens:      0,
				NumWorkingPods: 4,
				NumBrokenPods:  1,
				NumTurns:       40,
				AliensRespawn:  false,
				ActiveCards: map[string]int{
					"Red Card":   24,
					"Green Card": 26,
					"White Card": 4,
					"Adrenaline": 3,
					"Attack":     1,
					"Cat":        2,
					"Clone":      1,
					"Defense":    1,
					"Mutation":   1,
					"Sedatives":  1,
					"Sensor":     1,
					"Spotlight":  2,
					"Teleport":   1,
				},
				ActiveRoles: map[string]int{
					"Captain":           0,
					"Pilot":             0,
					"Copilot":           0,
					"Soldier":           0,
					"Psychologist":      99,
					"Executive Officer": 0,
					"Medic":             0,
					"Engineer":          0,
					"Fast":              99,
					"Surge":             0,
					"Blink":             0,
					"Silent":            0,
					"Brute":             0,
					"Invisible":         0,
					"Lurking":           0,
					"Psychic":           0,
				},
				RequiredRoles: map[string]int{
					"Captain":           0,
					"Pilot":             0,
					"Copilot":           0,
					"Soldier":           0,
					"Psychologist":      0,
					"Executive Officer": 0,
					"Medic":             0,
					"Engineer":          0,
					"Fast":              0,
					"Surge":             0,
					"Blink":             0,
					"Silent":            0,
					"Brute":             0,
					"Invisible":         0,
					"Lurking":           0,
					"Psychic":           0,
				},
			}),
		},
	}
}

func GetConfigAsJsonString(gameConfig GameConfig) string {
	if asJson, err := json.Marshal(gameConfig); err == nil {
		return string(asJson)
	}
	return "Error"
}
