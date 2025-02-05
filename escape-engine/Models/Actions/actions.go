package Actions

import (
	"encoding/json"
	"escape-engine/Models"
	"fmt"
	"log"
	"math"
	"math/rand"
	"slices"
)

const (
	Action_Attack   = "Attack"
	Action_Movement = "Movement"
	Action_Noise    = "Noise"
	Action_EndTurn  = "EndTurn"
	Action_PlayCard = "PlayCard"
)

// This is the way the frontend will send data to the backend during gameplay. They will
// send one of these objects, then the Rule Engine will take it, perform any updates to the
// internal model of the Game, then respond with a Changelog
type SubmittedAction struct {
	//The type of Turn you want to take. Should match exactly with the name of one of the below structs (i.e. "Movement", "Insertion", etc)
	Type string `json:"type"`
	//The actual turn object. Should have all the fields within the struct that you're wanting
	Turn json.RawMessage `json:"turn"`
	//ID of the player who is trying to submit this action. This is now supplied by the backend
	PlayerId string `json:"playerId"`
}

type Movement struct {
	ToRow string `json:"toRow"`
	ToCol int    `json:"toCol"`
}

// Set Row & Col to !, -99 to indicate no attack
type Attack struct {
	Row string `json:"row"`
	Col int    `json:"col"`
}

func (attack Attack) IsAttacking() bool {
	return attack.Row != "!" && attack.Col != -99
}

type Noise struct {
	Row string `json:"row"`
	Col int    `json:"col"`
}

func (noise Noise) IsNoisy() bool {
	return noise.Row != "!" && noise.Col != -99
}

type EndTurn struct {
}

type PlayCard struct {
	Name string `json:"name"`
}

func (move Movement) Execute(gameState *Models.GameState, playerId string) (Models.MovementEvent, error) {
	movement := Models.MovementEvent{}
	// if gameState.CurrentPlayer != playerId {
	// 	return nil, fmt.Errorf("player trying to execute turn is not current player")
	// }

	actingPlayer := gameState.GetCurrentPlayer()

	//Bounds check
	spaceKey := fmt.Sprintf("%s-%d", move.ToRow, move.ToCol)

	if space, exists := gameState.GameMap.Spaces[spaceKey]; exists {
		//Make sure it's an open space
		cantMoveInto := []int{
			Models.Space_AlienStart,
			Models.Space_HumanStart,
			Models.Space_Wall,
			Models.Space_UsedPod,
		}
		if slices.ContainsFunc(cantMoveInto, func(t int) bool { return space.Type == t }) {
			return movement, fmt.Errorf("space [%s,%d] is not allowed to be moved into", move.ToRow, move.ToCol)
		}

		//Make sure it's close enough
		// allowedSpaces := 1 //TODO: Figure out how to deal with aliens getting 3 spaces of movement after a kill
		// if actingPlayer.Team == PlayerTeam_Alien {
		// 	allowedSpaces = 2
		// }
		// if !checkMovement(move.ToRow, actingPlayer.Row, move.ToCol, actingPlayer.Col, allowedSpaces) {
		// 	return fmt.Errorf("movement not allowed") TODO: Turned off for now because hex grids make counting spaces HARD
		// }

		//At this point, player is allowed to execute the move
		actingPlayer.Row, actingPlayer.Col = move.ToRow, move.ToCol
		movement.NewRow, movement.NewCol = actingPlayer.Row, actingPlayer.Col
	} else {
		return movement, fmt.Errorf("requested space [%s-%d] not found in map", move.ToRow, move.ToCol)
	}

	return movement, nil
}

func checkMovement(toRow int, fromRow int, toCol int, fromCol int, allowedMovement int) bool {
	log.Printf("Row (%d->%d) Col(%d->%d), %d", fromRow, toRow, fromCol, toCol, allowedMovement)

	if int(math.Abs(float64(toRow-fromRow))) > allowedMovement {
		return false
	}

	allowedMovementOffset := allowedMovement - 1
	if fromCol%2 == 0 { //Even column
		log.Println("Even column")
		if toRow-fromRow > 0 {
			log.Println("Moving down")
			return int(math.Abs(float64(toCol-fromCol))) <= allowedMovementOffset
		} else {
			log.Println("moving up")
			return int(math.Abs(float64(toCol-fromCol))) <= allowedMovement
		}
	} else { //Odd column
		log.Println("Odd column")
		if toRow-fromRow < 0 { //Moving up
			log.Println("Moving up")
			return int(math.Abs(float64(toCol-fromCol))) <= allowedMovementOffset
		} else {
			log.Println("moving down")
			return int(math.Abs(float64(toCol-fromCol))) <= allowedMovement
		}
	}
}

func (attack Attack) Execute(gameState *Models.GameState, playerId string) (*Models.GameEvent, error) {
	actingPlayer := gameState.GetCurrentPlayer()

	var gameEvent *Models.GameEvent = &Models.GameEvent{
		Row:         attack.Row,
		Col:         attack.Col,
		Description: fmt.Sprintf("%s attacked [%s-%d]!\n", actingPlayer.Name, attack.Row, attack.Col),
	}
	alienStarts := gameState.GameMap.GetSpacesOfType(Models.Space_AlienStart)

	for index, player := range gameState.Players {
		if player.Id == actingPlayer.Id { //Don't kill the player doing the attacking
			continue
		}
		if player.Row != attack.Row || player.Col != attack.Col {
			continue
		}

		newSpaceForPlayer := alienStarts[rand.Intn(len(alienStarts))]

		gameState.Players[index].Team = Models.PlayerTeam_Alien
		gameState.Players[index].Row, gameState.Players[index].Col = newSpaceForPlayer.Row, newSpaceForPlayer.Col
		gameEvent.Description = string(fmt.Appendf([]byte(gameEvent.Description), "%s died!\n", gameState.Players[index].Name))
	}

	return gameEvent, nil
}

func DrawCard(gameState *Models.GameState, playerId string) (Models.CardEvent, error) {
	event := Models.CardEvent{}

	actingPlayer := gameState.GetCurrentPlayer()

	currentSpace := Models.Space{
		Row: actingPlayer.Row,
		Col: actingPlayer.Col,
	}

	if gameState.GameMap.Spaces[currentSpace.GetMapKey()].Type == Models.Space_Safe || gameState.GameMap.Spaces[currentSpace.GetMapKey()].Type == Models.Space_Pod {
		event.Type = Models.Card_NoCard
	} else { //TODO: Full implementation
		drawnCard := *drawRandomCardFromDeck(gameState)
		event.Card = drawnCard
		event.Type = drawnCard.GetType()
		if drawnCard.GetType() == Models.Card_White {
			actingPlayer.Hand.Cards = append(actingPlayer.Hand.Cards, drawnCard)
		} else {
			gameState.DiscardPile.Cards = append(gameState.DiscardPile.Cards, drawnCard)
		}
	}

	return event, nil
}

func (noise Noise) Execute(gameState *Models.GameState, playerId string) (*Models.GameEvent, error) {
	actingPlayer := gameState.GetCurrentPlayer()

	event := Models.GameEvent{
		Row: noise.Row,
		Col: noise.Col,
	}
	if noise.IsNoisy() {
		event.Description = fmt.Sprintf("Player '%s' made noise at [%s-%d]!", actingPlayer.Name, noise.Row, noise.Col)
	} else {
		event.Description = fmt.Sprintf("Player '%s' avoided making noise", actingPlayer.Name)
	}
	return &event, nil
}

func (endTurn EndTurn) Execute(gameState *Models.GameState, playerId string) (*Models.GameState, *Models.GameEvent, error) {
	var gameEvent *Models.GameEvent = nil

	actingPlayerIndex := slices.IndexFunc(gameState.Players, func(p Models.Player) bool { return playerId == p.Id })
	if actingPlayerIndex == -1 {
		return gameState, nil, fmt.Errorf("could not find acting player with ID == {%s}", playerId)
	}

	actingPlayer := &(gameState.Players[actingPlayerIndex])

	space := gameState.GameMap.Spaces[Models.Space{Row: actingPlayer.Row, Col: actingPlayer.Col}.GetMapKey()]
	if space.Type == Models.Space_Pod {
		if actingPlayer.Team == Models.PlayerTeam_Alien {
			return gameState, nil, fmt.Errorf("aliens are not allowed to enter pods")
		}

		totalPodCards := gameState.GameConfig.NumWorkingPods + gameState.GameConfig.NumBrokenPods
		if totalPodCards == 0 {
			return gameState, nil, fmt.Errorf("no pod cards left")
		}
		podCard := (rand.Intn(totalPodCards) + 1)
		podIsWorking := podCard > gameState.GameConfig.NumWorkingPods
		if gameState.GameConfig.NumBrokenPods == 0 { //0 broken pods is an edge case that will effectively make 1 "working" card act as a broken card
			podIsWorking = true
		}
		if gameState.GameConfig.NumWorkingPods == 0 { //0 working pods is an edge case in the opposite direction
			podIsWorking = false
		}
		if podIsWorking {
			gameEvent = &Models.GameEvent{
				Row:         actingPlayer.Row,
				Col:         actingPlayer.Col,
				Description: fmt.Sprintf("Player %s escaped using the Pod at [%s-%d]!", actingPlayer.Name, actingPlayer.Row, actingPlayer.Col),
			}
			actingPlayer.Team = Models.PlayerTeam_Spectator
			actingPlayer.Row, actingPlayer.Col = "!", -99

			gameState.GameConfig.NumWorkingPods -= 1
			gameState.GameMap.Spaces[space.GetMapKey()] = Models.Space{
				Row:  space.Row,
				Col:  space.Col,
				Type: Models.Space_UsedPod,
			}

			return gameState, gameEvent, nil
		} else {
			gameEvent = &Models.GameEvent{
				Row:         actingPlayer.Row,
				Col:         actingPlayer.Col,
				Description: fmt.Sprintf("Player %s tried the Pod at [%s-%d], but it didn't work!", actingPlayer.Name, actingPlayer.Row, actingPlayer.Col),
			}

			gameState.GameConfig.NumBrokenPods -= 1
			gameState.GameMap.Spaces[space.GetMapKey()] = Models.Space{
				Row:  space.Row,
				Col:  space.Col,
				Type: Models.Space_UsedPod,
			}
		}
	}

	gameState.CurrentPlayer = getNextPlayerId(gameState.Players, actingPlayerIndex)
	return gameState, gameEvent, nil //TODO: End the game when no human players remain
}

func (play PlayCard) Execute(gameState *Models.GameState, playerId string) (Models.GameEvent, error) {
	actingPlayer := gameState.GetCurrentPlayer()

	playedCardIndex := slices.IndexFunc(actingPlayer.Hand.Cards, func(c Models.Card) bool { return c.GetName() == play.Name })
	if playedCardIndex == -1 {
		return Models.GameEvent{}, fmt.Errorf("could not find card '%s' in Player's hand", play.Name)
	}

	playedCard := &(actingPlayer.Hand.Cards[playedCardIndex])
	//Important. Copy the card since slices.DeleteFunc will 0 out that location in memory
	cardCopy := *playedCard
	//Remove the card from the player's hand
	actingPlayer.Hand.Cards = slices.DeleteFunc(actingPlayer.Hand.Cards, func(c Models.Card) bool {
		return c == cardCopy
	})

	cardCopy.Play(gameState)

	gameState.DiscardPile.Cards = append(gameState.DiscardPile.Cards, cardCopy)

	return Models.GameEvent{
		Description: fmt.Sprintf("Player '%s' used a '%s' card!", actingPlayer.Name, cardCopy.GetName()),
		Row:         "!",
		Col:         -99,
	}, nil
}

func getNextPlayerId(players []Models.Player, currentIndex int) string {
	if !slices.ContainsFunc(players, func(p Models.Player) bool { return p.Team != Models.PlayerTeam_Spectator }) {
		return ""
	}

	nextIndex := currentIndex + 1
	if currentIndex >= len(players)-1 {
		nextIndex = 0
	}

	if players[nextIndex].Team == Models.PlayerTeam_Spectator {
		return getNextPlayerId(players, nextIndex)
	} else {
		return players[nextIndex].Id
	}
}

func drawRandomCardFromDeck(gameState *Models.GameState) *Models.Card {

	//Auto reshuffle
	if len(gameState.Deck.Cards) <= 0 {
		gameState.Deck.Cards = gameState.DiscardPile.Cards
		gameState.DiscardPile.Cards = []Models.Card{}
	}

	drawnCard := &(gameState.Deck.Cards[rand.Intn(len(gameState.Deck.Cards))])
	//Important. Copy the card since slices.DeleteFunc will 0 out that location in memory
	cardCopy := *drawnCard
	//Remove the card from the deck
	gameState.Deck.Cards = slices.DeleteFunc(gameState.Deck.Cards, func(c Models.Card) bool {
		return c == cardCopy
	})
	//Return address to copy of deleted card
	return &cardCopy
}
