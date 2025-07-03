package Models

import (
	"escape-engine/Models/Recap"
	"fmt"
	"math/rand"
	"slices"
)

const (
	Card_Red    = "Red"
	Card_Green  = "Green"
	Card_White  = "White"
	Card_NoCard = "SilentSector"
)

type CardBase struct {
	//Name of the Card
	Name string `json:"name"`
	//A Description of what this card does
	Description string `json:"description"`
	//One of the Card_* constants
	Type string `json:"type"`
	//A boolean describing whether this card should be completely discarded after a player plays it (true) or if it can go back into the discard pile (false)
	DestroyOnUse bool `json:"destroyOnUse"`
}

func (r *CardBase) GetName() string {
	return r.Name
}

func (r *CardBase) GetType() string {
	return r.Type
}

func (r *CardBase) GetDescription() string {
	return r.Description
}

func (r *CardBase) GetDestroyOnUse() bool {
	return r.DestroyOnUse
}

func (r *CardBase) SetDestroyOnUse(destroyOnUse bool) {
	r.DestroyOnUse = destroyOnUse
}

type CardPlayDetails struct {
	TargetRow    int
	TargetCol    string
	TargetPlayer string
}

// #region Red Card

type RedCard struct {
	CardBase
}

func (r *RedCard) Play(gameState *GameState, details CardPlayDetails) string {
	return ""
}

func NewRedCard() *RedCard {
	return &RedCard{
		CardBase: CardBase{
			Name:        "Red Card",
			Description: "Make a noise in the sector you just moved into",
			Type:        Card_Red,
		},
	}
}

// #region Green Card

type GreenCard struct {
	CardBase
}

func (g *GreenCard) Play(gameState *GameState, details CardPlayDetails) string {
	return ""
}

func NewGreenCard() *GreenCard {
	return &GreenCard{
		CardBase: CardBase{
			Name:        "Green Card",
			Description: "Make a noise in any sector of your choosing",
			Type:        Card_Green,
		},
	}
}

// #region White Card (No Item)

type WhiteCard struct {
	CardBase
}

func (w *WhiteCard) Play(gameState *GameState, details CardPlayDetails) string {
	return ""
}

func NewWhiteCard() *WhiteCard {
	return &WhiteCard{
		CardBase: CardBase{
			Name:        "White Card",
			Description: "You make no noise in this sector",
			Type:        Card_White,
		},
	}
}

// #region Adrenaline

type Adrenaline struct {
	CardBase
}

func (a *Adrenaline) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s.Name == StatusEffect_AdrenalineSurge }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, NewAdrenalineSurge())
	}
	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s", a.Name))
	return fmt.Sprintf("Player '%s' played an Adrenaline card. They can move one space farther on their next turn!", activePlayer.Name)
}

func NewAdrenaline() *Adrenaline {
	return &Adrenaline{
		CardBase: CardBase{
			Name:        "Adrenaline",
			Description: "Gives the Adrenaline Rush Status Effect, allowing one extra space of movement on your next movement",
			Type:        Card_White,
		},
	}
}

// #region Mutation

type Mutation struct {
	CardBase
}

func (m *Mutation) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	activePlayer.Team = PlayerTeam_Alien

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s", m.Name))
	return fmt.Sprintf("Player '%s' used a Mutation card! They've turned into an Alien!", activePlayer.Name)
}

func NewMutation() *Mutation {
	return &Mutation{
		CardBase: CardBase{
			Name:        "Mutation",
			Description: "Turns you into an Alien!",
			Type:        Card_White,
		},
	}
}

// #region Teleport

type Teleport struct {
	CardBase
}

func (t Teleport) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	humanStarts := gameState.GameMap.GetSpacesOfType(Space_HumanStart)
	toMoveTo := humanStarts[rand.Intn(len(humanStarts))]

	if activePlayer.Team == PlayerTeam_Alien {
		alienStarts := gameState.GameMap.GetSpacesOfType(Space_AlienStart)
		toMoveTo = alienStarts[rand.Intn(len(alienStarts))]
	}

	activePlayer.Row, activePlayer.Col = toMoveTo.Row, toMoveTo.Col

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s and teleported to [%s-%d]", t.Name, activePlayer.Col, activePlayer.Row))

	return fmt.Sprintf("Player '%s' used a Teleport card! They've been moved to a random starting sector of their team!", activePlayer.Name)
}

func NewTeleport() *Teleport {
	return &Teleport{
		CardBase: CardBase{
			Name:        "Teleport",
			Description: "Teleports you to a random Start Sector belonging to your team",
			Type:        Card_White,
		},
	}
}

// #region Clone

type Clone struct {
	CardBase
}

func (c Clone) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s.Name == StatusEffect_Cloned }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, NewCloned())
	}

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s", c.Name))

	return fmt.Sprintf("Player '%s' used a Clone card! They now have a clone ready in case they die!", activePlayer.Name)
}

func NewClone() *Clone {
	return &Clone{
		CardBase: CardBase{
			Name:        "Clone",
			Description: "Gives the Cloned Status Effect, creating a Clone of you that automatically activates upon death",
			Type:        Card_White,
		},
	}
}

// #region Defence

type Defense struct {
	CardBase
}

func (c Defense) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s.Name == StatusEffect_Armored }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, NewArmored())
	}

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s", c.Name))

	return fmt.Sprintf("Player '%s' used a Defense card! They'll be protected from the next attack!", activePlayer.Name)
}

func NewDefense() *Defense {
	return &Defense{
		CardBase: CardBase{
			Name:        "Defense",
			Description: "Gives the Armored Status Effect, protecting you from the next attack to hit you",
			Type:        Card_White,
		},
	}
}

// #region Spotlight

type Spotlight struct {
	CardBase
}

func (c Spotlight) Play(gameState *GameState, details CardPlayDetails) string {
	if details.TargetRow == -99 || details.TargetCol == "" {
		panic("No details provided for spotlight")
	}

	activePlayer := gameState.GetCurrentPlayer()

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played a %s on [%s-%d]", c.Name, details.TargetCol, details.TargetRow))

	descriptionString := fmt.Sprintf("Player '%s' used a Spotlight on [%s-%d]!", activePlayer.Name, details.TargetCol, details.TargetRow)

	seenPlayers := []Player{}

	adjacentSpaces := gameState.GameMap.GetSpacesWithinNthAdjacency(1, GetMapKey(details.TargetRow, details.TargetCol), nil)

	for spaceKey := range adjacentSpaces {
		for i, player := range gameState.Players {
			playerSpace := GetMapKey(player.Row, player.Col)
			if spaceKey == playerSpace {
				if gameState.Players[i].SubtractStatusEffect(StatusEffect_Invisible) {
					continue
				}
				seenPlayers = append(seenPlayers, player)
			}
		}
	}

	//Check space spotlight was played on as well
	for i, player := range gameState.Players {
		playerSpace := GetMapKey(player.Row, player.Col)
		spaceKey := GetMapKey(details.TargetRow, details.TargetCol)
		if spaceKey == playerSpace {
			if gameState.Players[i].SubtractStatusEffect(StatusEffect_Invisible) {
				continue
			}
			seenPlayers = append(seenPlayers, player)
		}
	}

	if len(seenPlayers) > 0 {
		for _, player := range seenPlayers {
			descriptionString += fmt.Sprintf("\nPlayer '%s' was seen at %s!", player.Name, GetMapKey(player.Row, player.Col))
		}
	} else {
		descriptionString += " No players were found!"
	}

	return descriptionString
}

func NewSpotlight() *Spotlight {
	return &Spotlight{
		CardBase: CardBase{
			Name:        "Spotlight",
			Description: "Reveals any player in the chosen space or any adjacent space",
			Type:        Card_White,
		},
	}
}

// #region Attack

type AttackCard struct {
	CardBase
}

func (c AttackCard) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s Card", c.Name))

	descriptionString := fmt.Sprintf("Player '%s' used an Attack Card! ", activePlayer.Name)

	gameEvent, _ := AttackSpace(activePlayer.Row, activePlayer.Col, gameState)

	descriptionString += gameEvent.Description

	return descriptionString
}

func NewAttackCard() *AttackCard {
	return &AttackCard{
		CardBase: CardBase{
			Name:        "Attack",
			Description: "Attacks the space you are currently in",
			Type:        Card_White,
		},
	}
}

// #region Sedatives

type Sedatives struct {
	CardBase
}

func (c Sedatives) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s.Name == StatusEffect_Sedated }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, NewSedated())
	}

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s", c.Name))

	return fmt.Sprintf("Player '%s' used Sedatives! They will treat the next sector they enter as a safe sector!", activePlayer.Name)
}

func NewSedatives() *Sedatives {
	return &Sedatives{
		CardBase: CardBase{
			Name:        "Sedatives",
			Description: "Gives the Sedated Status Effect, causing you to treat the next sector you enter as a Safe Sector",
			Type:        Card_White,
		},
	}
}

// #region Sensor

type Sensor struct {
	CardBase
}

func (c Sensor) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()
	indexOfTargetedPlayer := slices.IndexFunc(gameState.Players, func(p Player) bool { return p.Id == details.TargetPlayer })
	targetedPlayer := gameState.Players[indexOfTargetedPlayer]

	if gameState.Players[indexOfTargetedPlayer].SubtractStatusEffect(StatusEffect_Invisible) {
		go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s on Player '%s', but they were Invisible", c.Name, targetedPlayer.Name))
		return fmt.Sprintf("Player '%s' used a Sensor on Player '%s' but Player '%s' is Invisible!", activePlayer.Name, targetedPlayer.Name, targetedPlayer.Name)
	}

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s on Player '%s', revealing their location at [%s-%d]", c.Name, targetedPlayer.Name, targetedPlayer.Col, targetedPlayer.Row))

	return fmt.Sprintf("Player '%s' used a Sensor on Player '%s'! Player '%s' is at [%s-%d]", activePlayer.Name, targetedPlayer.Name, targetedPlayer.Name, targetedPlayer.Col, targetedPlayer.Row)
}

func NewSensor() *Sensor {
	return &Sensor{
		CardBase: CardBase{
			Name:        "Sensor",
			Description: "Publicly reveals the exact location of a player of your choice",
			Type:        Card_White,
		},
	}
}

// #region Cat

type Cat struct {
	CardBase
}

func (c Cat) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s.Name == StatusEffect_Feline }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, NewFeline())
	}
	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s", c.Name))
	return fmt.Sprintf("Player '%s' used a Cat! They can make 1 extra noise the next time they make a noise!", activePlayer.Name)
}

func NewCat() *Cat {
	return &Cat{
		CardBase: CardBase{
			Name:        "Cat",
			Description: "Gives the Feline StatusEffect, allowing you to pick one extra space to make a noise in, the next time you make a noise",
			Type:        Card_White,
		},
	}
}

// #region Scanner

type Scanner struct {
	CardBase
}

func (s Scanner) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()
	indexOfTargetedPlayer := slices.IndexFunc(gameState.Players, func(p Player) bool { return p.Id == details.TargetPlayer })
	targetedPlayer := gameState.Players[indexOfTargetedPlayer]

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s on Player '%s', revealing they were a %s %s", s.Name, targetedPlayer.Name, targetedPlayer.Role, targetedPlayer.Team))

	return fmt.Sprintf("Player '%s' used a Scanner on Player '%s'! Player '%s' is a %s %s!", activePlayer.Name, targetedPlayer.Name, targetedPlayer.Name, targetedPlayer.Role, targetedPlayer.Team)
}

func NewScanner() *Scanner {
	return &Scanner{
		CardBase: CardBase{
			Name:        "Scanner",
			Description: "Publicly reveals the Team & Role of a Player of your choice",
			Type:        Card_White,
		},
	}
}
