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
	StatusEffect_Deceptive       = "Deceptive"
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
		Description: "You may move 1 extra space on your next movement",
		UsesLeft:    1,
	}
}

// #region Cloned

func NewCloned() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Cloned,
		Description: "You have a clone that will automatically activate upon death",
		UsesLeft:    1,
	}
}

// #region Armored

func NewArmored() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Armored,
		Description: "You are defended from the next attack that hits you",
		UsesLeft:    1,
	}
}

// #region Hyperphagic

func NewHyperphagic() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Hyperphagic,
		Description: "You have fed on a human, gaining strength. But you want more...",
		UsesLeft:    1,
	}
}

// #region Sedated

func NewSedated() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Sedated,
		Description: "You are sedated and will treat the next space you enter as a Safe Sector",
		UsesLeft:    1,
	}
}

// #region Feline

func NewFeline() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Feline,
		Description: "You can choose 1 additional Sector to make noise in, the next time they make any noise",
		UsesLeft:    1,
	}
}

// #region Invisible

func NewInvisible() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Invisible,
		Description: "You are immune to Spotlights and Sensors",
		UsesLeft:    1,
	}
}

// #region Lurking
func NewLurking() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Lurking,
		Description: "You can end your turn on the same space you started on",
		UsesLeft:    1,
	}
}

// #region Knowhow
func NewKnowhow() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Knowhow,
		Description: "You draw 2 Escape Pod cards upon arriving at an escape pod",
		UsesLeft:    1,
	}
}

// #region Deceptive
func NewDeceptive() StatusEffect {
	return StatusEffect{
		Name:        StatusEffect_Deceptive,
		Description: "You may make a noise in any sector upon drawing a White Card, as if you had drawn a Green Card instead",
		UsesLeft:    1,
	}
}
