package Models

import (
	"fmt"
	"math/rand"
	"slices"
)

func GetAllowedSpaces(player *Player, gameState *GameState) int {
	allowedSpaces := 1

	//Adrenaline
	if indexOfEffect := slices.IndexFunc(player.StatusEffects, func(s StatusEffect) bool { return s.GetName() == StatusEffect_AdrenalineSurge }); indexOfEffect != -1 {
		player.StatusEffects[indexOfEffect].SubtractUse(player)
		allowedSpaces++
	}

	//Aliens
	if player.Team == PlayerTeam_Alien {
		allowedSpaces++
	}

	//Hyperphagic
	if indexOfEffect := slices.IndexFunc(player.StatusEffects, func(s StatusEffect) bool { return s.GetName() == StatusEffect_Hyperphagic }); indexOfEffect != -1 {
		allowedSpaces++
		//Hyperphagic is a permanent bonus, so don't subtract any uses
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

func AttackSpace(row string, col int, gameState GameState) (*GameEvent, error) {
	actingPlayer := gameState.GetCurrentPlayer()

	var gameEvent *GameEvent = &GameEvent{
		Row:         row,
		Col:         col,
		Description: fmt.Sprintf("Player %s attacked [%s-%d]!\n", actingPlayer.Name, row, col),
	}
	alienStarts := gameState.GameMap.GetSpacesOfType(Space_AlienStart)

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
			return gameState.StatusEffectPriorities[s2] - gameState.StatusEffectPriorities[s1]
		})

		playerWasSaved := false
		for _, se := range defenseEffects {
			if player.HasStatusEffect(se) {
				switch se {
				case StatusEffect_Armored:
					playerWasSaved = true
					gameEvent.Description += fmt.Sprintf("Player %s was saved by Armor!\n", player.Name)
				case StatusEffect_Cloned:
					playerWasSaved = true
					gameEvent.Description += fmt.Sprintf("Player %s activated their Emergency Clone!\n", player.Name)
				}
			}
			if playerWasSaved {
				break
			}
		}

		if !playerWasSaved {
			newSpaceForPlayer := alienStarts[rand.Intn(len(alienStarts))]

			gameState.Players[index].Team = PlayerTeam_Alien
			gameState.Players[index].Row, gameState.Players[index].Col = newSpaceForPlayer.Row, newSpaceForPlayer.Col
			gameEvent.Description += fmt.Sprintf("%s died!\n", player.Name)
		}
	}

	return gameEvent, nil
}

func GetUnmarshalledCardArray(intermediate []struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}) []Card {

	cards := make([]Card, len(intermediate))

	for i, card := range intermediate {
		switch card.Name {
		case "Red Card":
			cards[i] = NewRedCard()
		case "Green Card":
			cards[i] = NewGreenCard()
		case "Mutation":
			cards[i] = NewMutation()
		case "Adrenaline":
			cards[i] = NewAdrenaline()
		case "Teleport":
			cards[i] = NewTeleport()
		case "Clone":
			cards[i] = NewClone()
		case "Defence":
			cards[i] = NewDefense()
		case "Spotlight":
			cards[i] = NewSpotlight()
		case "Attack":
			cards[i] = NewAttackCard()
		case "Sedatives":
			cards[i] = NewSedatives()
		}
	}

	return cards
}
