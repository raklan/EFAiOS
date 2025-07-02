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

	//Hyperphagic
	if player.HasStatusEffect(StatusEffect_Hyperphagic) {
		//Hyperphagic is a permanent bonus, so don't subtract any uses
		allowedSpaces++
	}

	return allowedSpaces
}

func GetNonmovableSpaces(player *Player) []int {
	cantMoveInto := []int{
		Space_AlienStart,
		Space_HumanStart,
		Space_Wall,
		Space_UsedPod,
	}
	if player.Team == PlayerTeam_Alien {
		cantMoveInto = append(cantMoveInto, Space_Pod)
	}

	return cantMoveInto
}

func AttackSpace(row int, col string, gameState GameState) (*GameEvent, error) {
	actingPlayer := gameState.GetCurrentPlayer()

	var gameEvent *GameEvent = &GameEvent{
		Row:         row,
		Col:         col,
		Description: fmt.Sprintf("Player '%s' attacked [%s-%d]!\n", actingPlayer.Name, col, row),
	}
	alienStarts := gameState.GameMap.GetSpacesOfType(Space_AlienStart)

	go Recap.AddDataToRecap(gameState.Id, actingPlayer.Id, gameState.Turn, fmt.Sprintf("Attacked [%s-%d]", col, row))

	for index, player := range gameState.Players {
		if player.Id == actingPlayer.Id { //Don't kill the player doing the attacking
			continue
		}
		if player.Row != row || player.Col != col {
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
			if player.SubtractStatusEffect(se) {
				switch se {
				case StatusEffect_Armored:
					playerWasSaved = true
					gameEvent.Description += fmt.Sprintf("Player '%s' was saved by Armor!\n", player.Name)
					go Recap.AddDataToRecap(gameState.Id, actingPlayer.Id, gameState.Turn, fmt.Sprintf("Attacked Player '%s'", player.Name))
					go Recap.AddDataToRecap(gameState.Id, player.Id, gameState.Turn, fmt.Sprintf("Attacked by Player '%s' and saved by Armor", actingPlayer.Name))
				case StatusEffect_Cloned:
					playerWasSaved = true
					gameEvent.Description += fmt.Sprintf("Player '%s' activated their Emergency Clone!\n", player.Name)
					go Recap.AddDataToRecap(gameState.Id, actingPlayer.Id, gameState.Turn, fmt.Sprintf("Attacked Player '%s'", player.Name))
					go Recap.AddDataToRecap(gameState.Id, player.Id, gameState.Turn, fmt.Sprintf("Killed by player '%s' and activated Emergency Clone", actingPlayer.Name))
				}
			}
			if playerWasSaved {
				break
			}
		}

		if !playerWasSaved {
			if player.Team == PlayerTeam_Human || (player.Team == PlayerTeam_Alien && gameState.GameMap.GameConfig.AliensRespawn) {
				newSpaceForPlayer := alienStarts[rand.Intn(len(alienStarts))]

				gameState.Players[index].Team = PlayerTeam_Alien
				gameState.Players[index].Row, gameState.Players[index].Col = newSpaceForPlayer.Row, newSpaceForPlayer.Col
			} else {
				gameState.Players[index].Team = PlayerTeam_Spectator
				gameState.Players[index].Row, gameState.Players[index].Col = -99, "!"
			}

			gameEvent.Description += fmt.Sprintf("Player '%s' died!\n", player.Name)
			go Recap.AddDataToRecap(gameState.Id, actingPlayer.Id, gameState.Turn, fmt.Sprintf("Killed Player '%s'", player.Name))
			go Recap.AddDataToRecap(gameState.Id, player.Id, gameState.Turn, fmt.Sprintf("Killed by Player '%s'", actingPlayer.Name))
		}
	}

	return gameEvent, nil
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
