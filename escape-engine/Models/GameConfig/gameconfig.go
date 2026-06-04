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
	//A map of Player ID -> Team name. A player whose ID does not appear here will be randomly assigned a team.
	TeamAssignments map[string]string `json:"teamAssignments"`
	//Number of Working Escape Pods left. The Game automatically ends when this number hits 0.
	NumWorkingPods int `json:"numWorkingPods"`
	//Number of Broken Escape Pods left
	NumBrokenPods int `json:"numBrokenPods"`
	//Which cards should be active, as well as how many of each
	ActiveCards map[string]int `json:"activeCards"`
	//Which roles should be active, as well as the maximum number allowed to be present. Should be >= that role's presence in RequiredRoles, if it's required
	ActiveRoles map[string]int `json:"activeRoles"`
	//Which roles should be guaranteed to be in the game, as well as the number of players that should have that role
	RequiredRoles map[string]int `json:"requiredRoles"`
	//All modifiers for this game
	Modifiers GameModifiers `json:"modifiers"`
	//Which StatusEffects should be active, as well as their priority
	ActiveStatusEffects map[string]int `json:"activeStatusEffects"`
}

type GameModifiers struct {
	//Number of turns before the game should end
	NumTurns int `json:"numTurns"`
	//Whether players' turns should automatically end after moving (and choosing to attack, for aliens) if there are no cards in their hand
	AutoTurnEnd bool `json:"autoTurnEnd"`
	//Whether the game should spawn a random Evacuation Sector on turn # in `EvacuationTiming`
	EvacuationMode bool `json:"evacuationMode"`
	//On this turn, spawn a randomly-placed Evacuation Sector
	EvacuationTiming int `json:"evacuationTiming"`
	//Whether Aliens should spawn in a randomly selected Escape Pod, instead of an Alien Start Sector.
	InfestedPodsMode bool `json:"infestedPodsMode"`
	//Whether to block all escape pods until there is only one human remaining
	LastManStandingMode bool `json:"lastManStandingMode"`
	//Whether the game should add a randomly placed guaranteed escape sector after all escape pods are marked as used.
	LastResortMode bool `json:"lastResortMode"`
	//Whether the game should be played such that escape pods do not remove humans from the game, and the humans' new win condition is to activate each escape pod
	ReactorMode bool `json:"reactorMode"`
	//Whether Alien players should respawn in an Alien Start Sector upon death, instead of joining the spectator team.
	RelentlessAliensMode bool `json:"relentlessAliensMode"`
	//Whether players should be placed in a completely random sector to start the game, instead of their team's Start Sectors.
	ScatterMode bool `json:"scatterMode"`
	//Whether the game should be played such that after NumTurns, the Humans automatically win, instead of automatically dying
	SurvivalMode bool `json:"survivalMode"`
	//Whether the Escape Pods should start the game as blocked, and become unblocked using the logic in PodUnblockTiming
	UnstablePodsMode bool `json:"unstablePodsMode"`
	//When to mark Escape Pods as unblocked. The following logic applies: If == -2, unblock on Even Turns. If == -1, unblock on Odd turns. If == 0, never unblock. If > 0, unblock forever starting on that turn #
	PodUnblockTiming int `json:"podUnblockTiming"`
}

//Descriptions that can be more or less plugged directly into an innerHTML attribute of the modifier descriptions window. Modes that have configuration values have a [%VAR] in the string that can be replaced as needed.
var FormattedModifierDescriptions = map[string]string{
	"evacuationMode":       "<span class=\"modifier-entry-title\">Evacuation Mode</span>: At the beginning of Turn [%VAR], a random Dangerous or Safe Sector will become an Evacuation Sector. Any human reaching this sector escapes.",
	"infestedPodsMode":     "<span class=\"modifier-entry-title\">Infested Pods Mode</span>: Aliens start in a randomly selected Escape Pod Sector.",
	"lastManStandingMode":  "<span class=\"modifier-entry-title\">Last Man Standing Mode</span>: All Escape Pod Sectors are unusable until there is only 1 Human Player remaining.",
	"lastResortMode":       "<span class=\"modifier-entry-title\">Last Resort Mode</span>: Once all Escape Pod Sectors are used, a random Dangerous or Safe Sector will become an Evacuation Sector. Any human reaching this sector escapes.",
	"reactorMode":          "<span class=\"modifier-entry-title\">Reactor Mode</span>: Humans cannot escape through Escape Pod Sectors. When all Escape Pod sectors have been visited, all remaining Humans automatically escape.",
	"relentlessAliensMode": "<span class=\"modifier-entry-title\">Relentless Aliens Mode</span>: Upon death, Alien Players will respawn at the start of their next turn.",
	"scatterMode":          "<span class=\"modifier-entry-title\">Scatter Mode</span>: All players start in a randomly selected Sector.",
	"survivalMode":         "<span class=\"modifier-entry-title\">Survival Mode</span>: After [%VAR] turns have elapsed, all surviving Humans automatically escape.",
	"unstablePodsMode":     "<span class=\"modifier-entry-title\">Unstable Pods Mode</span>: All Escape Pod Sectors are unusable. They only become usable [%VAR].",
}

func GetConfigPresets() []GameConfigPreset {
	return []GameConfigPreset{
		{
			Name:        "Tabletop",
			Description: "The Classic experience, with all settings exactly as they appear in the Tabletop version.",
			ConfigJson: GetConfigAsJsonString(GameConfig{
				NumHumans:      0,
				NumAliens:      0,
				NumWorkingPods: 4,
				NumBrokenPods:  1,
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
				Modifiers: GameModifiers{
					NumTurns:             40,
					RelentlessAliensMode: false,
					AutoTurnEnd:          false,
				},
			}),
		},
		{
			Name:        "The Basics",
			Description: "A simplified version of the Tabletop preset, with all Roles turned off and every Item card replaced with a White \"Silent\" Card. Ideal for people just learning the game.",
			ConfigJson: GetConfigAsJsonString(GameConfig{
				NumHumans:      0,
				NumAliens:      0,
				NumWorkingPods: 4,
				NumBrokenPods:  1,
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
				Modifiers: GameModifiers{
					NumTurns:             40,
					RelentlessAliensMode: false,
					AutoTurnEnd:          false,
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
				Modifiers: GameModifiers{
					NumTurns:             40,
					RelentlessAliensMode: false,
					AutoTurnEnd:          false,
				},
			}),
		},
		{
			Name:        "Empty",
			Description: "Everything set to 0. Tailor the game exactly how you like it.",
			ConfigJson: GetConfigAsJsonString(GameConfig{
				NumHumans:      0,
				NumAliens:      0,
				NumWorkingPods: 0,
				NumBrokenPods:  0,
				ActiveCards: map[string]int{
					"Red Card":   0,
					"Green Card": 0,
					"White Card": 0,
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
				Modifiers: GameModifiers{
					NumTurns:             0,
					RelentlessAliensMode: false,
					AutoTurnEnd:          false,
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
