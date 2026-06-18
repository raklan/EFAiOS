package Models

import (
	"escape-engine/Models/Recap"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

func GetAllowedSpaces(player *Player, gameState *GameState) int {
	log.Println("getting allowed spaces for player", player)

	allowedSpaces := 1

	//Adrenaline
	if player.SubtractStatusEffect(StatusEffect_AdrenalineSurge) {
		allowedSpaces++
	}

	//Aliens
	if player.Team == PlayerTeam_Alien {
		allowedSpaces++
	}

	//Hyperphagic is a permanent bonus, so don't subtract any uses
	if player.HasStatusEffect(StatusEffect_Hyperphagic) {
		if gameState.GameMap.GameConfig.Modifiers.BloodlustMode {
			for _, se := range player.StatusEffects {
				if se.Name == StatusEffect_Hyperphagic {
					allowedSpaces += se.UsesLeft
				}
			}
		} else {
			allowedSpaces++
		}
	}

	return allowedSpaces
}

func GetNonmovableSpaces(player *Player) []int {
	cantMoveInto := []int{
		Space_AlienStart,
		Space_HumanStart,
		Space_Wall,
		Space_UsedPod,
		Space_BlockedPod,
	}
	if player.Team == PlayerTeam_Alien {
		cantMoveInto = append(cantMoveInto, Space_Pod)
	}

	return cantMoveInto
}

func AttackSpace(row int, col string, gameState *GameState) (*GameEvent, error) {
	actingPlayer := gameState.GetCurrentPlayer()

	var gameEvent *GameEvent = &GameEvent{
		Row:         row,
		Col:         col,
		Description: fmt.Sprintf("Player '%s' attacked [%s-%d]!\n", actingPlayer.Name, col, row),
	}

	potentialNewStarts := []Space{}
	if gameState.GameMap.GameConfig.Modifiers.InfestedPodsMode {
		potentialNewStarts = gameState.GameMap.GetSpacesOfType(Space_Pod, Space_UsedPod)
	} else if gameState.GameMap.GameConfig.Modifiers.ScatterMode {
		potentialNewStarts = gameState.GameMap.GetSpacesOfType(Space_AlienStart, Space_HumanStart, Space_Dangerous, Space_Safe)
	} else {
		potentialNewStarts = gameState.GameMap.GetSpacesOfType(Space_AlienStart)
	}

	go Recap.AddDataToRecap(gameState.Id, actingPlayer.Id, gameState.Turn, fmt.Sprintf("Attacked [%s-%d]", col, row))

	for index, targetPlayer := range gameState.Players {
		if targetPlayer.Id == actingPlayer.Id { //Don't kill the player doing the attacking
			continue
		}
		if targetPlayer.Row != row || targetPlayer.Col != col {
			continue
		}

		//Check if the player has any status effects that will save them, apply in order of priority
		defenseEffects := []string{
			StatusEffect_Armored,
			StatusEffect_Cloned,
		}

		slices.SortFunc(defenseEffects, func(s1 string, s2 string) int {
			return gameState.GameMap.GameConfig.ActiveStatusEffects[s2] - gameState.GameMap.GameConfig.ActiveStatusEffects[s1]
		})

		playerWasSaved := false
		for _, se := range defenseEffects {
			if targetPlayer.SubtractStatusEffect(se) {
				switch se {
				case StatusEffect_Armored:
					playerWasSaved = true
					if !gameState.GameMap.GameConfig.Modifiers.UnconfirmedKillsMode {
						gameEvent.Description += fmt.Sprintf("Player '%s' was saved by Armor!\n", targetPlayer.Name)
					}
					go Recap.AddDataToRecap(gameState.Id, actingPlayer.Id, gameState.Turn, fmt.Sprintf("Attacked Player '%s'", targetPlayer.Name))
					go Recap.AddDataToRecap(gameState.Id, targetPlayer.Id, gameState.Turn, fmt.Sprintf("Attacked by Player '%s' and saved by Armor", actingPlayer.Name))
				case StatusEffect_Cloned:
					playerWasSaved = true
					if !gameState.GameMap.GameConfig.Modifiers.UnconfirmedKillsMode {
						gameEvent.Description += fmt.Sprintf("Player '%s' activated their Emergency Clone!\n", targetPlayer.Name)
					}

					humanStarts := gameState.GameMap.GetSpacesOfType(Space_HumanStart)
					newSpaceForPlayer := humanStarts[rand.Intn(len(humanStarts))]
					gameState.Players[index].Row, gameState.Players[index].Col = newSpaceForPlayer.Row, newSpaceForPlayer.Col

					go Recap.AddDataToRecap(gameState.Id, actingPlayer.Id, gameState.Turn, fmt.Sprintf("Attacked Player '%s'", targetPlayer.Name))
					go Recap.AddDataToRecap(gameState.Id, targetPlayer.Id, gameState.Turn, fmt.Sprintf("Killed by player '%s' and activated Emergency Clone", actingPlayer.Name))
				}
			}
			if playerWasSaved {
				break
			}
		}

		if !playerWasSaved {
			if targetPlayer.Team == PlayerTeam_Human || (targetPlayer.Team == PlayerTeam_Alien && gameState.GameMap.GameConfig.Modifiers.RelentlessAliensMode) {
				newSpaceForPlayer := potentialNewStarts[rand.Intn(len(potentialNewStarts))]

				gameState.Players[index].Team = PlayerTeam_Alien
				gameState.Players[index].Row, gameState.Players[index].Col = newSpaceForPlayer.Row, newSpaceForPlayer.Col
			} else {
				gameState.Players[index].Team = PlayerTeam_Spectator
				gameState.Players[index].Row, gameState.Players[index].Col = -99, "!"
			}

			for _, card := range targetPlayer.Hand {
				if !card.GetDestroyOnUse() {
					gameState.DiscardPile = append(gameState.DiscardPile, card)
				}
			}
			gameState.Players[index].Hand = []Card{}
			gameState.Players[index].StatusEffects = []StatusEffect{}

			if gameState.GameMap.GameConfig.Modifiers.NecrophagiaMode {
				actingPlayer.AddStatusEffect(StatusEffect_Hyperphagic, NewHyperphagic)
			} else {
				if actingPlayer.Team == PlayerTeam_Alien && targetPlayer.Team == PlayerTeam_Human {
					actingPlayer.AddStatusEffect(StatusEffect_Hyperphagic, NewHyperphagic)
				}
			}

			if !gameState.GameMap.GameConfig.Modifiers.UnconfirmedKillsMode {
				gameEvent.Description += fmt.Sprintf("Player '%s' died!\n", targetPlayer.Name)
			}
			go Recap.AddDataToRecap(gameState.Id, actingPlayer.Id, gameState.Turn, fmt.Sprintf("Killed Player '%s'", targetPlayer.Name))
			go Recap.AddDataToRecap(gameState.Id, targetPlayer.Id, gameState.Turn, fmt.Sprintf("Killed by Player '%s'", actingPlayer.Name))
		}
	}

	if gameState.GameMap.GameConfig.Modifiers.DescructiveAttacksMode {
		attackedSpace := gameState.GameMap.Spaces[Space{
			Row: row,
			Col: col,
		}.GetMapKey()]

		newSpaceType := Space_Wall
		if attackedSpace.Type == Space_Safe {
			newSpaceType = Space_Dangerous
		}

		gameState.GameMap.Spaces[attackedSpace.GetMapKey()] = Space{
			Row:  attackedSpace.Row,
			Col:  attackedSpace.Col,
			Type: newSpaceType,
		}
	}

	return gameEvent, nil
}

func HandleEscapePodBlocking(gameState *GameState) {
	podsShouldBeBlocked := false
	if gameState.GameMap.GameConfig.Modifiers.UnstablePodsMode {
		//Logic is recorded in docs for PodUnblockTiming
		switch gameState.GameMap.GameConfig.Modifiers.PodUnblockTiming {
		case -2:
			podsShouldBeBlocked = gameState.Turn%2 != 0
		case -1:
			podsShouldBeBlocked = gameState.Turn%2 == 0
		case 0:
			podsShouldBeBlocked = true
		default:
			podsShouldBeBlocked = gameState.Turn < gameState.GameMap.GameConfig.Modifiers.PodUnblockTiming
		}
	}

	if gameState.GameMap.GameConfig.Modifiers.LastManStandingMode {
		podsShouldBeBlocked = podsShouldBeBlocked || len(gameState.GetHumanPlayers()) > 1
	}

	//Don't touch used pods
	allPods := gameState.GameMap.GetSpacesOfType(Space_Pod, Space_BlockedPod)
	typeToSet := Space_Pod
	for _, pod := range allPods {
		if podsShouldBeBlocked {
			typeToSet = Space_BlockedPod
		}

		gameState.GameMap.Spaces[pod.GetMapKey()] = Space{
			Row:  pod.Row,
			Col:  pod.Col,
			Type: typeToSet,
		}
	}

}

func GetUnmarshalledCardArray(intermediate []CardBase) []Card {

	cards := make([]Card, len(intermediate))

	for i, card := range intermediate {
		switch card.Name {
		case "Red Card":
			cards[i] = NewRedCard()
		case "Green Card":
			cards[i] = NewGreenCard()
		case "White Card":
			cards[i] = NewWhiteCard()
		case "Mutation":
			cards[i] = NewMutation()
		case "Adrenaline":
			cards[i] = NewAdrenaline()
		case "Teleport":
			cards[i] = NewTeleport()
		case "Clone":
			cards[i] = NewClone()
		case "Defense":
			cards[i] = NewDefense()
		case "Spotlight":
			cards[i] = NewSpotlight()
		case "Attack":
			cards[i] = NewAttackCard()
		case "Sedatives":
			cards[i] = NewSedatives()
		case "Sensor":
			cards[i] = NewSensor()
		case "Cat":
			cards[i] = NewCat()
		case "Scanner":
			cards[i] = NewScanner()
		case "Unstable Teleporter":
			cards[i] = NewUnstableTeleporter()
		case "Hiding Spot":
			cards[i] = NewHidingSpot()
		case "Cloaking Device":
			cards[i] = NewCloakingDevice()
		case "Engineering Manual":
			cards[i] = NewEngineeringManual()
		case "Noisemaker":
			cards[i] = NewNoisemaker()
		}
		cards[i].SetDestroyOnUse(card.DestroyOnUse)
	}

	return cards
}

// Gets a pseudo-random (Key, Value) pair from the given map. Runs in O(len(theMap)) time
func GetRandomMapPair[T1 comparable, T2 any](theMap map[T1]T2) (T1, T2) {
	// Yeah, this is pretty janky, but it's simple and I don't need speed for maps that are guaranteed to be small. Sue me
	randomIndex := rand.Intn(len(theMap))
	index := 0
	for k, v := range theMap {
		if index == randomIndex {
			return k, v
		}
		index++
	}
	panic("How did we get here?")
}

// Randomly generates a number in the range [0,n), ensuring the number is not equal to [exclude]
func RandExclusive(n int, exclude int) int {
	//Janky but good enough here
	numPicked := rand.Intn(n)
	for numPicked == exclude {
		numPicked = rand.Intn(n)
	}
	return numPicked
}
