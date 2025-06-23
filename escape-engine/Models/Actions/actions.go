package Actions

import (
	"encoding/json"
	"escape-engine/Models"
	"fmt"
	"maps"
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
	ToRow int    `json:"toRow"`
	ToCol string `json:"toCol"`
}

// Set Row & Col to !, -99 to indicate no attack
type Attack struct {
	Row int    `json:"row"`
	Col string `json:"col"`
}

func (attack Attack) IsAttacking() bool {
	return attack.Row != -99 && attack.Col != "!"
}

type Noise struct {
	Row  int    `json:"row"`
	Col  string `json:"col"`
	Row2 int    `json:"row2"`
	Col2 string `json:"col2"`
}

func (noise Noise) IsNoisy() bool {
	return noise.Row != -99 && noise.Col != "!"
}

func (noise Noise) IsNoisy2() bool {
	return noise.Row2 != -99 && noise.Col2 != "!"
}

type EndTurn struct {
}

type PlayCard struct {
	Name         string `json:"name"`
	Row          int    `json:"row"`
	Col          string `json:"col"`
	TargetPlayer string `json:"targetPlayer"`
}

func (move Movement) Execute(gameState *Models.GameState, playerId string) (Models.MovementEvent, error) {
	movement := Models.MovementEvent{}
	// if gameState.CurrentPlayer != playerId {
	// 	return nil, fmt.Errorf("player trying to execute turn is not current player")
	// }

	actingPlayer := gameState.GetCurrentPlayer()

	//Bounds check
	spaceKey := Models.GetMapKey(move.ToRow, move.ToCol)

	if space, exists := gameState.GameMap.Spaces[spaceKey]; exists {
		//Make sure it's an open space

		if slices.ContainsFunc(Models.GetNonmovableSpaces(actingPlayer), func(t int) bool { return space.Type == t }) {
			return movement, fmt.Errorf("space [%s] is not allowed to be moved into", spaceKey)
		}

		//Make sure it's close enough
		numAllowedSpaces := Models.GetAllowedSpaces(actingPlayer, gameState)
		allowedSpacesToMoveTo := gameState.GameMap.GetSpacesWithinNthAdjacency(numAllowedSpaces, actingPlayer.GetSpaceMapKey(), actingPlayer)

		if _, exists := allowedSpacesToMoveTo[space.GetMapKey()]; !exists {
			if actingPlayer.GetSpaceMapKey() == space.GetMapKey() {
				if !actingPlayer.SubtractStatusEffect(Models.StatusEffect_Lurking) {
					return movement, fmt.Errorf("you cannot end your turn on the same space you started on")
				}
				//If they're not moving but have Lurking, it's ok
			} else {
				return movement, fmt.Errorf("space [%s] is too far away", spaceKey)
			}
		}

		//At this point, player is allowed to execute the move
		actingPlayer.Row, actingPlayer.Col = move.ToRow, move.ToCol
		movement.NewRow, movement.NewCol = actingPlayer.Row, actingPlayer.Col
	} else {
		return movement, fmt.Errorf("requested space [%s-%d] not found in map", move.ToCol, move.ToRow)
	}

	return movement, nil
}

func GetPotentialMoves(gameState *Models.GameState, playerId string) []string {
	actingPlayer := gameState.GetCurrentPlayer()
	numAllowedSpaces := Models.GetAllowedSpaces(actingPlayer, gameState)
	allowedSpacesToMoveTo := gameState.GameMap.GetSpacesWithinNthAdjacency(numAllowedSpaces, actingPlayer.GetSpaceMapKey(), actingPlayer)
	if actingPlayer.HasStatusEffect(Models.StatusEffect_Lurking) {
		allowedSpacesToMoveTo[actingPlayer.GetSpaceMapKey()] = gameState.GameMap.Spaces[actingPlayer.GetSpaceMapKey()]
	}

	return slices.Collect(maps.Keys(allowedSpacesToMoveTo))
}

func (attack Attack) Execute(gameState *Models.GameState, playerId string) (*Models.GameEvent, error) {
	return Models.AttackSpace(attack.Row, attack.Col, *gameState)
}

func DrawCard(gameState *Models.GameState, playerId string) (Models.CardEvent, error) {
	event := Models.CardEvent{}

	actingPlayer := gameState.GetCurrentPlayer()

	currentSpace := Models.Space{
		Row: actingPlayer.Row,
		Col: actingPlayer.Col,
	}

	if actingPlayer.SubtractStatusEffect(Models.StatusEffect_Sedated) {
		event.Type = Models.Card_NoCard
	} else if gameState.GameMap.Spaces[currentSpace.GetMapKey()].Type == Models.Space_Safe ||
		gameState.GameMap.Spaces[currentSpace.GetMapKey()].Type == Models.Space_Pod {
		event.Type = Models.Card_NoCard
	} else {
		drawnCard := *drawRandomCardFromDeck(gameState)
		event.Card = drawnCard
		event.Type = drawnCard.GetType()
		if drawnCard.GetType() == Models.Card_White {
			if actingPlayer.Team == Models.PlayerTeam_Human && drawnCard.GetName() != "White Card" { //May need tweaking. Currently discards item cards picked up by Aliens
				actingPlayer.Hand = append(actingPlayer.Hand, drawnCard)
			}

			if actingPlayer.SubtractStatusEffect(Models.StatusEffect_Deceptive) {
				event.Type = Models.Card_Green
				gameState.DiscardPile = append(gameState.DiscardPile, drawnCard)
			}
		} else {
			//Don't need to check if the card needs to be destroyed in this case because it'll only ever be Red/Green/NoItem cards at this point
			gameState.DiscardPile = append(gameState.DiscardPile, drawnCard)
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
		if noise.IsNoisy2() {
			if !actingPlayer.SubtractStatusEffect(Models.StatusEffect_Feline) {
				return &Models.GameEvent{}, fmt.Errorf("player '%s' does not have 'Feline' StatusEffect", actingPlayer.Name)
			}

			//Randomize which space appears first and which appears second as an extra layer of secrecy
			if rand.Intn(11)%2 == 0 {
				event.Description = fmt.Sprintf("Player '%s' made noise at [%s-%d] and [%s-%d]!", actingPlayer.Name, noise.Col, noise.Row, noise.Col2, noise.Row2)
			} else {
				event.Description = fmt.Sprintf("Player '%s' made noise at [%s-%d] and [%s-%d]!", actingPlayer.Name, noise.Col2, noise.Row2, noise.Col, noise.Row)
			}
		} else {
			event.Description = fmt.Sprintf("Player '%s' made noise at [%s-%d]!", actingPlayer.Name, noise.Col, noise.Row)
		}
	} else {
		event.Description = fmt.Sprintf("Player '%s' avoided making noise", actingPlayer.Name)
	}
	return &event, nil
}

func (endTurn EndTurn) Execute(gameState *Models.GameState, playerId string) (*Models.GameEvent, error) {
	var gameEvent *Models.GameEvent = nil

	actingPlayer := gameState.GetCurrentPlayer()

	space := gameState.GameMap.Spaces[actingPlayer.GetSpaceMapKey()]
	if space.Type == Models.Space_Pod {
		if actingPlayer.Team == Models.PlayerTeam_Alien {
			return nil, fmt.Errorf("aliens are not allowed to enter pods")
		}

		totalPodCards := gameState.GameMap.GameConfig.NumWorkingPods + gameState.GameMap.GameConfig.NumBrokenPods
		if totalPodCards == 0 {
			return nil, fmt.Errorf("no pod cards left")
		}
		workingPods := slices.Repeat([]bool{true}, gameState.GameMap.GameConfig.NumWorkingPods)
		brokenPods := slices.Repeat([]bool{false}, gameState.GameMap.GameConfig.NumBrokenPods)

		allPods := slices.Concat(workingPods, brokenPods)
		drawnPod := rand.Intn(len(allPods))

		podIsWorking := allPods[drawnPod]
		if podIsWorking {
			gameEvent = &Models.GameEvent{
				Row:         actingPlayer.Row,
				Col:         actingPlayer.Col,
				Description: fmt.Sprintf("Player '%s' escaped using the Pod at [%s-%d]!", actingPlayer.Name, actingPlayer.Col, actingPlayer.Row),
			}
			actingPlayer.Team = Models.PlayerTeam_Spectator
			actingPlayer.Row, actingPlayer.Col = -99, "!"

			gameState.GameMap.GameConfig.NumWorkingPods -= 1

		} else {
			if len(allPods) > 1 && actingPlayer.SubtractStatusEffect(Models.StatusEffect_Knowhow) {
				secondDrawnPod := Models.RandExclusive(len(allPods), drawnPod)
				if allPods[secondDrawnPod] {
					gameEvent = &Models.GameEvent{
						Row:         actingPlayer.Row,
						Col:         actingPlayer.Col,
						Description: fmt.Sprintf("Player '%s' escaped using the Pod at [%s-%d]!", actingPlayer.Name, actingPlayer.Col, actingPlayer.Row),
					}
					actingPlayer.Team = Models.PlayerTeam_Spectator
					actingPlayer.Row, actingPlayer.Col = -99, "!"

					gameState.GameMap.GameConfig.NumWorkingPods -= 1

				}
			} else {
				gameEvent = &Models.GameEvent{
					Row:         actingPlayer.Row,
					Col:         actingPlayer.Col,
					Description: fmt.Sprintf("Player '%s' tried the Pod at [%s-%d], but it didn't work!", actingPlayer.Name, actingPlayer.Col, actingPlayer.Row),
				}

				gameState.GameMap.GameConfig.NumBrokenPods -= 1
			}
		}
		gameState.GameMap.Spaces[space.GetMapKey()] = Models.Space{
			Row:  space.Row,
			Col:  space.Col,
			Type: Models.Space_UsedPod,
		}
	}

	gameState.CurrentPlayer = getNextPlayerId(gameState.Players, slices.IndexFunc(gameState.Players, func(p Models.Player) bool { return p.Id == actingPlayer.Id }))

	if gameState.CurrentPlayer == gameState.Players[0].Id {
		gameState.Turn++
	}

	if gameState.Turn > gameState.GameMap.GameConfig.NumTurns {
		if gameEvent == nil {
			gameEvent = &Models.GameEvent{
				Row:         -99,
				Col:         "!",
				Description: "",
			}
		}

		gameEvent.Description += "The ship ran out of oxygen!"

		for i, player := range gameState.Players {
			if player.Team == Models.PlayerTeam_Human {
				gameState.Players[i].Team = Models.PlayerTeam_Spectator
				gameEvent.Description += fmt.Sprintf(" Player '%s' died!", player.Name)
			}
		}
	}

	return gameEvent, nil
}

func (play PlayCard) Execute(gameState *Models.GameState, playerId string) (Models.GameEvent, error) {
	actingPlayer := gameState.GetCurrentPlayer()

	playedCardIndex := slices.IndexFunc(actingPlayer.Hand, func(c Models.Card) bool { return c.GetName() == play.Name })
	if playedCardIndex == -1 {
		return Models.GameEvent{}, fmt.Errorf("could not find card '%s' in Player's hand", play.Name)
	}

	playedCard := &(actingPlayer.Hand[playedCardIndex])
	//Important. Copy the card since slices.DeleteFunc will 0 out that location in memory
	cardCopy := *playedCard
	//Remove the card from the player's hand
	actingPlayer.Hand = slices.DeleteFunc(actingPlayer.Hand, func(c Models.Card) bool {
		return c == cardCopy
	})

	cardMessage := cardCopy.Play(gameState, Models.CardPlayDetails{
		TargetRow:    play.Row,
		TargetCol:    play.Col,
		TargetPlayer: play.TargetPlayer,
	})

	if !cardCopy.GetDestroyOnUse() {
		gameState.DiscardPile = append(gameState.DiscardPile, cardCopy)
	}

	return Models.GameEvent{
		Description: cardMessage,
		Row:         -99,
		Col:         "!",
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
	if len(gameState.Deck) <= 0 {
		gameState.Deck = gameState.DiscardPile
		gameState.DiscardPile = []Models.Card{}
	}

	drawnCard := &(gameState.Deck[rand.Intn(len(gameState.Deck))])
	//Important. Copy the card since slices.DeleteFunc will 0 out that location in memory
	cardCopy := *drawnCard
	//Remove the card from the deck
	gameState.Deck = slices.DeleteFunc(gameState.Deck, func(c Models.Card) bool {
		return c == cardCopy
	})
	//Return address to copy of deleted card
	return &cardCopy
}
