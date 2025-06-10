package Models

var RoleAssigners = map[string]func(*Player){
	Role_Captain: AssignCaptain,
	Role_Pilot:   AssignPilot,
	Role_Soldier: AssignSoldier,

	Role_SpeedyAlien: AssignSpeedyAlien,
	Role_BlinkAlien:  AssignBlinkAlien,
	Role_SilentAlien: AssignSilentAlien,
}

var RoleTeams = map[string]string{
	Role_Captain: PlayerTeam_Human,
	Role_Pilot:   PlayerTeam_Human,
	Role_Soldier: PlayerTeam_Human,

	Role_SpeedyAlien: PlayerTeam_Alien,
	Role_BlinkAlien:  PlayerTeam_Alien,
	Role_SilentAlien: PlayerTeam_Alien,
}

const (
	Role_Captain = "Captain"
	Role_Pilot   = "Pilot"
	Role_Soldier = "Soldier"

	Role_SpeedyAlien = "Speedy Alien"
	Role_BlinkAlien  = "Blink Alien"
	Role_SilentAlien = "Silent Alien"
)

//#region Human Roles

func AssignCaptain(player *Player) {
	player.Role = Role_Captain
	player.StatusEffects = append(player.StatusEffects, NewSedated())
}

func AssignPilot(player *Player) {
	player.Role = Role_Pilot
	player.Hand = append(player.Hand, NewCat())
}

func AssignSoldier(player *Player) {
	player.Role = Role_Soldier
	player.Hand = append(player.Hand, NewAttackCard())
}

//#region Alien Roles

func AssignSpeedyAlien(player *Player) {
	player.Role = Role_SpeedyAlien
	player.StatusEffects = append(player.StatusEffects, NewAdrenalineSurge())
}

func AssignBlinkAlien(player *Player) { //TODO: This can allow the blink alien to instantly get to the human sector
	player.Role = Role_BlinkAlien
	player.Hand = append(player.Hand, NewTeleport())
}

func AssignSilentAlien(player *Player) {
	player.Role = Role_SilentAlien
	player.Hand = append(player.Hand, NewSedatives())
}
