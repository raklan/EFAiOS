package Models

import (
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
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type CardPlayDetails struct {
	TargetRow    string
	TargetCol    int
	TargetPlayer string
}

// #region Red Card

type RedCard struct {
	CardBase
}

func (r *RedCard) GetName() string {
	return r.Name
}

func (r *RedCard) GetType() string {
	return r.Type
}

func (r *RedCard) GetDescription() string {
	return r.Description
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

func (g *GreenCard) GetName() string {
	return g.Name
}

func (g *GreenCard) GetType() string {
	return g.Type
}

func (g *GreenCard) GetDescription() string {
	return g.Description
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

func (w *WhiteCard) GetName() string {
	return w.Name
}

func (w *WhiteCard) GetType() string {
	return w.Type
}

func (w *WhiteCard) GetDescription() string {
	return w.Description
}

func (w *WhiteCard) Play(gameState *GameState, details CardPlayDetails) string {
	return ""
}

func NewWhiteCard() *WhiteCard {
	return &WhiteCard{
		CardBase: CardBase{
			Name:        "White Card",
			Description: "You make no noise in this sector",
			Type:        Card_NoCard,
		},
	}
}

// #region Adrenaline

type Adrenaline struct {
	CardBase
}

func (a *Adrenaline) GetName() string {
	return a.Name
}

func (a *Adrenaline) GetType() string {
	return a.Type
}

func (a *Adrenaline) GetDescription() string {
	return a.Description
}

func (a *Adrenaline) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s.GetName() == StatusEffect_AdrenalineSurge }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, NewAdrenalineSurge())
	}

	return fmt.Sprintf("Player %s played an Adrenaline card. They can move one space farther on their next turn!", activePlayer.Name)
}

func NewAdrenaline() *Adrenaline {
	return &Adrenaline{
		CardBase: CardBase{
			Name:        "Adrenaline",
			Description: "Gives a rush of adrenaline, allowing one extra space of movement",
			Type:        Card_White,
		},
	}
}

// #region Mutation

type Mutation struct {
	CardBase
}

func (m *Mutation) GetName() string {
	return m.Name
}

func (m *Mutation) GetType() string {
	return m.Type
}

func (m *Mutation) GetDescription() string {
	return m.Description
}

func (m *Mutation) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	activePlayer.Team = PlayerTeam_Alien

	return fmt.Sprintf("Player %s used a Mutation card! They've turned into an Alien!", activePlayer.Name)
}

func NewMutation() *Mutation {
	return &Mutation{
		CardBase: CardBase{
			Name:        "Mutation",
			Description: "Turns the Player into an Alien!",
			Type:        Card_White,
		},
	}
}

// #region Teleport

type Teleport struct {
	CardBase
}

func (t Teleport) GetName() string {
	return t.Name
}

func (t Teleport) GetType() string {
	return t.Type
}

func (t Teleport) GetDescription() string {
	return t.Description
}

func (t Teleport) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	humanStarts := gameState.GameMap.GetSpacesOfType(Space_HumanStart)

	toMoveTo := humanStarts[rand.Intn(len(humanStarts))]

	activePlayer.Row, activePlayer.Col = toMoveTo.Row, toMoveTo.Col

	return fmt.Sprintf("Player %s used a Teleport card! They've been moved to a random Human starting sector!", activePlayer.Name)
}

func NewTeleport() *Teleport {
	return &Teleport{
		CardBase: CardBase{
			Name:        "Teleport",
			Description: "Teleports the Player to a random Human Start Sector",
			Type:        Card_White,
		},
	}
}

// #region Clone

type Clone struct {
	CardBase
}

func (c Clone) GetName() string {
	return c.Name
}

func (c Clone) GetType() string {
	return c.Type
}

func (c Clone) GetDescription() string {
	return c.Description
}

func (c Clone) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s.GetName() == StatusEffect_Cloned }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, NewCloned())
	}

	return fmt.Sprintf("Player %s used a Clone card! They now have a clone ready in case they die!", activePlayer.Name)
}

func NewClone() *Clone {
	return &Clone{
		CardBase: CardBase{
			Name:        "Clone",
			Description: "Creates a Clone of this player that automatically activates upon death",
			Type:        Card_White,
		},
	}
}

// #region Defence

type Defense struct {
	CardBase
}

func (c Defense) GetName() string {
	return c.Name
}

func (c Defense) GetType() string {
	return c.Type
}

func (c Defense) GetDescription() string {
	return c.Description
}

func (c Defense) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s.GetName() == StatusEffect_Armored }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, NewArmored())
	}

	return fmt.Sprintf("Player %s used a Defense card! They'll be protected from the next attack!", activePlayer.Name)
}

func NewDefense() *Defense {
	return &Defense{
		CardBase: CardBase{
			Name:        "Defense",
			Description: "Makes this player invulnerable to the next attack that hits them",
			Type:        Card_White,
		},
	}
}

// #region Spotlight

type Spotlight struct {
	CardBase
}

func (c Spotlight) GetName() string {
	return c.Name
}

func (c Spotlight) GetType() string {
	return c.Type
}

func (c Spotlight) GetDescription() string {
	return c.Description
}

func (c Spotlight) Play(gameState *GameState, details CardPlayDetails) string {
	if details.TargetRow == "" || details.TargetCol == -99 {
		panic("No details provided for spotlight")
	}

	activePlayer := gameState.GetCurrentPlayer()

	descriptionString := fmt.Sprintf("Player %s used a Spotlight!", activePlayer.Name)

	seenPlayers := []Player{}

	adjacentSpaces := gameState.GameMap.GetSpacesWithinNthAdjacency(1, GetMapKey(details.TargetRow, details.TargetCol))

	for spaceKey := range adjacentSpaces {
		for _, player := range gameState.Players {
			playerSpace := GetMapKey(player.Row, player.Col)
			if spaceKey == playerSpace {
				seenPlayers = append(seenPlayers, player)
			}
		}
	}

	//Check space spotlight was played on as well
	for _, player := range gameState.Players {
		playerSpace := GetMapKey(player.Row, player.Col)
		spaceKey := GetMapKey(details.TargetRow, details.TargetCol)
		if spaceKey == playerSpace {
			seenPlayers = append(seenPlayers, player)
		}
	}

	if len(seenPlayers) > 0 {
		for _, player := range seenPlayers {
			descriptionString += fmt.Sprintf("\nPlayer %s was seen at %s!", player.Name, GetMapKey(player.Row, player.Col))
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

func (c AttackCard) GetName() string {
	return c.Name
}

func (c AttackCard) GetType() string {
	return c.Type
}

func (c AttackCard) GetDescription() string {
	return c.Description
}

func (c AttackCard) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	descriptionString := fmt.Sprintf("Player %s used an Attack Card!", activePlayer.Name)

	gameEvent, _ := AttackSpace(activePlayer.Row, activePlayer.Col, *gameState)

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

func (c Sedatives) GetName() string {
	return c.Name
}

func (c Sedatives) GetType() string {
	return c.Type
}

func (c Sedatives) GetDescription() string {
	return c.Description
}

func (c Sedatives) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s.GetName() == StatusEffect_Sedated }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, NewSedated())
	}

	return fmt.Sprintf("Player %s used Sedatives!", activePlayer.Name)
}

func NewSedatives() *Sedatives {
	return &Sedatives{
		CardBase: CardBase{
			Name:        "Sedatives",
			Description: "Sedates the player, causing them to treat the next sector they enter as a Safe Sector",
			Type:        Card_White,
		},
	}
}
