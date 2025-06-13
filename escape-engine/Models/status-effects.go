package Models

import (
	"slices"
)

const (
	StatusEffect_AdrenalineSurge = "Adrenaline Surge"
	StatusEffect_Cloned          = "Cloned"
	StatusEffect_Armored         = "Armored"
	StatusEffect_Hyperphagic     = "Hyperphagic"
	StatusEffect_Sedated         = "Sedated"
	StatusEffect_Feline          = "Feline"
	StatusEffect_Invisible       = "Invisible"
	StatusEffect_Lurking         = "Lurking"
	StatusEffect_Knowhow         = "Knowhow"
)

type StatusEffect struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UsesLeft    int    `json:"usesLeft"`
}

func (a *StatusEffect) AddUse() int {
	a.UsesLeft++
	return a.UsesLeft
}

func (a *StatusEffect) SubtractUse(player *Player) bool {
	a.UsesLeft--

	if a.UsesLeft <= 0 {
		player.StatusEffects = slices.DeleteFunc(player.StatusEffects, func(s StatusEffect) bool { return s.Name == a.Name })
		return false
	}

	return true
}

// #region Adrenaline Surge

func NewAdrenalineSurge() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_AdrenalineSurge,
		Description: "The affected player may move 1 extra space",
		UsesLeft:    1,
	}
}

// #region Cloned

func NewCloned() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Cloned,
		Description: "This player has a clone that will automatically activate upon death",
		UsesLeft:    1,
	}
}

// #region Armored

func NewArmored() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Armored,
		Description: "This player is defended from the next attack that hits them",
		UsesLeft:    1,
	}
}

// #region Hyperphagic

func NewHyperphagic() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Hyperphagic,
		Description: "This Alien has fed on a human, gaining strength. But they want more...",
		UsesLeft:    1,
	}
}

// #region Sedated

func NewSedated() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Sedated,
		Description: "This player is sedated and will treat the next space they enter as a Safe Sector",
		UsesLeft:    1,
	}
}

// #region Feline

func NewFeline() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Feline,
		Description: "This player can make 2 noises the next time they make a noise",
		UsesLeft:    1,
	}
}

// #region Invisible

func NewInvisible() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Invisible,
		Description: "This player is immune to Spotlights and Sensors",
		UsesLeft:    1,
	}
}

// #region Lurking
func NewLurking() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Lurking,
		Description: "This player can end their turn on the same space they started on",
		UsesLeft:    1,
	}
}

// #region Knowhow
func NewKnowhow() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Knowhow,
		Description: "This player draws 2 Escape Pod cards upon arriving at an escape pod",
		UsesLeft:    1,
	}
}
