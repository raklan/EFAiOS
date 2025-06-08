package Models

import (
	"slices"
)

type StatusEffectBase struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UsesLeft    int    `json:"usesLeft"`
}

// #region Adrenaline Surge

const (
	StatusEffect_AdrenalineSurge = "Adrenaline Surge"
	StatusEffect_Cloned          = "Cloned"
	StatusEffect_Armored         = "Armored"
	StatusEffect_Hyperphagic     = "Hyperphagic"
)

type AdrenalineSurge struct {
	StatusEffectBase
}

func NewAdrenalineSurge() *AdrenalineSurge {
	return &AdrenalineSurge{
		StatusEffectBase: StatusEffectBase{
			Name:        StatusEffect_AdrenalineSurge,
			Description: "The affected player may move 1 extra space",
			UsesLeft:    1,
		},
	}
}

func (a *AdrenalineSurge) GetName() string {
	return a.Name
}

func (a *AdrenalineSurge) GetDescription() string {
	return a.Description
}

func (a *AdrenalineSurge) GetUsesLeft() int {
	return a.UsesLeft
}

func (a *AdrenalineSurge) AddUse() int {
	a.UsesLeft++
	return a.GetUsesLeft()
}

func (a *AdrenalineSurge) SubtractUse(player *Player) bool {
	a.UsesLeft--

	if a.UsesLeft <= 0 {
		player.StatusEffects = slices.DeleteFunc(player.StatusEffects, func(s StatusEffect) bool { return s == a })
		return false
	}

	return true
}

// #region Cloned

type Cloned struct {
	StatusEffectBase
}

func NewCloned() *Cloned {
	return &Cloned{
		StatusEffectBase: StatusEffectBase{
			Name:        StatusEffect_Cloned,
			Description: "This player has a clone that will automatically activate upon death",
			UsesLeft:    1,
		},
	}
}

func (s *Cloned) GetName() string {
	return s.Name
}

func (s *Cloned) GetDescription() string {
	return s.Description
}

func (s *Cloned) GetUsesLeft() int {
	return s.UsesLeft
}

func (s *Cloned) AddUse() int {
	s.UsesLeft++
	return s.GetUsesLeft()
}

func (s *Cloned) SubtractUse(player *Player) bool {
	s.UsesLeft--

	if s.UsesLeft <= 0 {
		player.StatusEffects = slices.DeleteFunc(player.StatusEffects, func(s2 StatusEffect) bool { return s2 == s })
		return false
	}

	return true
}

// #region Armored

type Armored struct {
	StatusEffectBase
}

func NewArmored() *Armored {
	return &Armored{
		StatusEffectBase: StatusEffectBase{
			Name:        StatusEffect_Armored,
			Description: "This player is defended from the next attack that hits them",
			UsesLeft:    1,
		},
	}
}

func (s *Armored) GetName() string {
	return s.Name
}

func (s *Armored) GetDescription() string {
	return s.Description
}

func (s *Armored) GetUsesLeft() int {
	return s.UsesLeft
}

func (s *Armored) AddUse() int {
	s.UsesLeft++
	return s.GetUsesLeft()
}

func (s *Armored) SubtractUse(player *Player) bool {
	s.UsesLeft--

	if s.UsesLeft <= 0 {
		player.StatusEffects = slices.DeleteFunc(player.StatusEffects, func(s2 StatusEffect) bool { return s2 == s })
		return false
	}

	return true
}

// #region Hyperphagic

type Hyperphagic struct {
	StatusEffectBase
}

func NewHyperphagic() *Armored {
	return &Armored{
		StatusEffectBase: StatusEffectBase{
			Name:        StatusEffect_Hyperphagic,
			Description: "This Alien has fed on a human, gaining strength. But they want more...",
			UsesLeft:    1,
		},
	}
}

func (s *Hyperphagic) GetName() string {
	return s.Name
}

func (s *Hyperphagic) GetDescription() string {
	return s.Description
}

func (s *Hyperphagic) GetUsesLeft() int {
	return s.UsesLeft
}

func (s *Hyperphagic) AddUse() int {
	s.UsesLeft++
	return s.GetUsesLeft()
}

func (s *Hyperphagic) SubtractUse(player *Player) bool {
	s.UsesLeft--

	if s.UsesLeft <= 0 {
		player.StatusEffects = slices.DeleteFunc(player.StatusEffects, func(s2 StatusEffect) bool { return s2 == s })
		return false
	}

	return true
}
