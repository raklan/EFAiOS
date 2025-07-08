package Models

var RoleAssigners = map[string]func(*Player){
	//Vanilla
	Role_Captain:          AssignCaptain,
	Role_Pilot:            AssignPilot,
	Role_Copilot:          AssignCopilot,
	Role_Soldier:          AssignSoldier,
	Role_Psychologist:     AssignPsychologist,
	Role_ExecutiveOfficer: AssignExecutiveOfficer,
	Role_Medic:            AssignMedic,
	Role_Engineer:         AssignEngineer,

	Role_FastAlien:      AssignFastAlien,
	Role_SurgeAlien:     AssignSurgeAlien,
	Role_BlinkAlien:     AssignBlinkAlien,
	Role_SilentAlien:    AssignSilentAlien,
	Role_BruteAlien:     AssignBruteAlien,
	Role_InvisibleAlien: AssignInvisibleAlien,
	Role_LurkingAlien:   AssignLurkingAlien,
	Role_PsychicAlien:   AssignPsychicAlien,

	Role_Scout:                 AssignScout,
	Role_CommunicationsOfficer: AssignCommunicationsOfficer,

	Role_TrackerAlien: AssignTrackerAlien,
	Role_CallingAlien: AssignCallingAlien,
}

var RoleTeams = map[string]string{
	//Vanilla
	Role_Captain:          PlayerTeam_Human,
	Role_Pilot:            PlayerTeam_Human,
	Role_Copilot:          PlayerTeam_Human,
	Role_Soldier:          PlayerTeam_Human,
	Role_Psychologist:     PlayerTeam_Human,
	Role_ExecutiveOfficer: PlayerTeam_Human,
	Role_Medic:            PlayerTeam_Human,
	Role_Engineer:         PlayerTeam_Human,

	Role_FastAlien:      PlayerTeam_Alien,
	Role_SurgeAlien:     PlayerTeam_Alien,
	Role_BlinkAlien:     PlayerTeam_Alien,
	Role_SilentAlien:    PlayerTeam_Alien,
	Role_BruteAlien:     PlayerTeam_Alien,
	Role_InvisibleAlien: PlayerTeam_Alien,
	Role_LurkingAlien:   PlayerTeam_Alien,
	Role_PsychicAlien:   PlayerTeam_Alien,

	//Raklan's Arsenal
	Role_Scout:                 PlayerTeam_Human,
	Role_CommunicationsOfficer: PlayerTeam_Human,

	Role_CallingAlien: PlayerTeam_Alien,
	Role_TrackerAlien: PlayerTeam_Alien,
}

var RoleDescriptions = map[string]string{
	//Vanilla
	Role_Captain:          "You start with 1 stack of the Sedated Status Effect",
	Role_Pilot:            "You start with 1 Cat card",
	Role_Copilot:          "You start with 1 Teleport card",
	Role_Soldier:          "You start with 1 Attack card",
	Role_Psychologist:     "You start the game in a randomly selected Alien start sector",
	Role_ExecutiveOfficer: "You start with 1 stack of the Lurking Status Effect",
	Role_Medic:            "You start with 1 Scanner card",
	Role_Engineer:         "You permanently gain the Knowhow Status Effect",

	Role_FastAlien:      "You start with 1 stack of the Adrenaline Surge Status Effect",
	Role_SurgeAlien:     "You start with 1 Adrenaline card",
	Role_BlinkAlien:     "You start with 1 Teleport card",
	Role_SilentAlien:    "You start with 1 Sedatives card",
	Role_BruteAlien:     "You permanently gain the Armored Status Effect",
	Role_InvisibleAlien: "You permanently gain the Invisible Status Effect",
	Role_LurkingAlien:   "You permanently gain the Lurking Status Effect",
	Role_PsychicAlien:   "You permanently gain the Deceptive Status Effect",

	//Raklan's Arsenal
	Role_Scout:                 "You start with 1 stack of the Invisible Status Effect",
	Role_CommunicationsOfficer: "You start with 1 Noisemaker card",

	Role_CallingAlien: "You start with 1 Cat card",
	Role_TrackerAlien: "You start with 1 Spotlight card",
}

const (
	//Vanilla
	Role_Captain          = "Captain"
	Role_Pilot            = "Pilot"
	Role_Copilot          = "Copilot"
	Role_Soldier          = "Soldier"
	Role_Psychologist     = "Psychologist"
	Role_ExecutiveOfficer = "Executive Officer"
	Role_Medic            = "Medic"
	Role_Engineer         = "Engineer"

	Role_FastAlien      = "Fast"
	Role_SurgeAlien     = "Surge"
	Role_BlinkAlien     = "Blink"
	Role_SilentAlien    = "Silent"
	Role_BruteAlien     = "Brute"
	Role_InvisibleAlien = "Invisible"
	Role_LurkingAlien   = "Lurking"
	Role_PsychicAlien   = "Psychic"

	//Raklan's Arsenal
	Role_Scout                 = "Scout"
	Role_CommunicationsOfficer = "Communications Officer"

	Role_CallingAlien = "Calling"
	Role_TrackerAlien = "Tracker"
)

//#region Human Roles

func AssignCaptain(player *Player) {
	player.Role = Role_Captain
	player.StatusEffects = append(player.StatusEffects, NewSedated())
}

func AssignPilot(player *Player) {
	player.Role = Role_Pilot
	cat := NewCat()
	cat.DestroyOnUse = true
	player.Hand = append(player.Hand, cat)
}

func AssignCopilot(player *Player) {
	player.Role = Role_Copilot
	tp := NewTeleport()
	tp.DestroyOnUse = true
	player.Hand = append(player.Hand, tp)
}

func AssignSoldier(player *Player) {
	player.Role = Role_Soldier
	attack := NewAttackCard()
	attack.DestroyOnUse = true
	player.Hand = append(player.Hand, attack)
}

func AssignPsychologist(player *Player) {
	player.Role = Role_Psychologist
}

func AssignExecutiveOfficer(player *Player) {
	player.Role = Role_ExecutiveOfficer
	player.StatusEffects = append(player.StatusEffects, NewLurking())
}

func AssignMedic(player *Player) {
	player.Role = Role_Medic
	scanner := NewScanner()
	scanner.DestroyOnUse = true
	player.Hand = append(player.Hand, scanner)
}

func AssignEngineer(player *Player) {
	player.Role = Role_Engineer
	knowhow := NewKnowhow()
	knowhow.UsesLeft = 1000
	player.StatusEffects = append(player.StatusEffects, knowhow)
}

//#region Alien Roles

func AssignFastAlien(player *Player) {
	player.Role = Role_FastAlien
	player.StatusEffects = append(player.StatusEffects, NewAdrenalineSurge())
}

func AssignSurgeAlien(player *Player) {
	player.Role = Role_SurgeAlien
	adr := NewAdrenaline()
	adr.DestroyOnUse = true
	player.Hand = append(player.Hand, adr)
}

func AssignBlinkAlien(player *Player) {
	player.Role = Role_BlinkAlien
	tp := NewTeleport()
	tp.DestroyOnUse = true
	player.Hand = append(player.Hand, tp)
}

func AssignSilentAlien(player *Player) {
	player.Role = Role_SilentAlien
	sed := NewSedatives()
	sed.DestroyOnUse = true
	player.Hand = append(player.Hand, sed)
}

func AssignBruteAlien(player *Player) {
	player.Role = Role_BruteAlien
	playerArmor := NewArmored()
	playerArmor.UsesLeft = 1000
	player.StatusEffects = append(player.StatusEffects, playerArmor)
}

func AssignInvisibleAlien(player *Player) {
	player.Role = Role_InvisibleAlien
	invis := NewInvisible()
	invis.UsesLeft = 1000
	player.StatusEffects = append(player.StatusEffects, invis)
}

func AssignLurkingAlien(player *Player) {
	player.Role = Role_LurkingAlien
	lurk := NewLurking()
	lurk.UsesLeft = 1000
	player.StatusEffects = append(player.StatusEffects, lurk)
}

func AssignPsychicAlien(player *Player) {
	player.Role = Role_PsychicAlien
	dec := NewDeceptive()
	dec.UsesLeft = 1000
	player.StatusEffects = append(player.StatusEffects, dec)
}

// #region Raklan's Arsenal
func AssignScout(player *Player) {
	player.Role = Role_Scout
	player.StatusEffects = append(player.StatusEffects, NewInvisible())
}

func AssignCommunicationsOfficer(player *Player) {
	player.Role = Role_CommunicationsOfficer
	nm := NewNoisemaker()
	nm.DestroyOnUse = true
	player.Hand = append(player.Hand, nm)
}

func AssignTrackerAlien(player *Player) {
	player.Role = Role_TrackerAlien
	spotlight := NewSpotlight()
	spotlight.DestroyOnUse = true
	player.Hand = append(player.Hand, spotlight)
}

func AssignCallingAlien(player *Player) {
	player.Role = Role_CallingAlien
	cat := NewCat()
	cat.DestroyOnUse = true
	player.Hand = append(player.Hand, cat)
}
