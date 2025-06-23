const GAME_CONFIG_DEFAULT = {
    workingPods: 4,
    brokenPods: 1,

    numTurns: 40,
    aliensRespawn: false,

    activeCards: {
        'Red Card': 24,
        'Green Card': 26,
        'White Card': 4,

        Adrenaline: 3,
        Attack: 1,
        Cat: 2,
        Clone: 1,
        Defense: 1,
        Mutation: 1,
        Sedatives: 1,
        Sensor: 1,
        Spotlight: 2,
        Teleport: 1
    },

    activeRoles: {
        Captain: 1,
        Pilot: 1,
        Copilot: 1,
        Soldier: 1,
        Psychologist: 1,
        'Executive Officer': 1,
        Medic: 1,
        Engineer: 1,

        Fast: 1,
        Surge: 1,
        Blink: 1,
        Silent: 1,
        Brute: 1,
        Invisible: 1,
        Lurking: 1, 
        Psychic: 1
    },
}

function setConfigFormFromObject(configObject) {
    setGeneralConfig(configObject);
    setCardConfig(configObject);
    setRoleConfig(configObject);
}

function setConfigFormFromString(configString){
    const configObject = JSON.parse(configString);
    setGeneralConfig(configObject);
    setCardConfig(configObject);
    setRoleConfig(configObject);
}

function setGeneralConfig(configObject) {
    let configForm = document.getElementById("lobby-gameConfig")

    configForm['config-numHumans'].value = configObject.numHumans;
    configForm['config-numAliens'].value = configObject.numAliens;

    configForm['config-numWorkingPods'].value = configObject.numWorkingPods;
    configForm['config-numBrokenPods'].value = configObject.numBrokenPods;

    configForm['config-numTurns'].value = configObject.numTurns;
    configForm['config-aliensRespawn'].checked = configObject.aliensRespawn
}

function setCardConfig(configObject) {
    let configForm = document.getElementById("lobby-gameConfig")

    configForm['config-numRedCards'].value = configObject.activeCards['Red Card'];
    configForm['config-numGreenCards'].value = configObject.activeCards['Green Card'];
    configForm['config-numWhiteCards'].value = configObject.activeCards['White Card'];

    configForm['config-numTeleport'].value = configObject.activeCards.Teleport;
    configForm['config-numClone'].value = configObject.activeCards.Clone;
    configForm['config-numDefense'].value = configObject.activeCards.Defense;

    configForm['config-numSpotlight'].value = configObject.activeCards.Spotlight;
    configForm['config-numAttack'].value = configObject.activeCards.Attack;
    configForm['config-numSensor'].value = configObject.activeCards.Sensor;

    configForm['config-numAdrenaline'].value = configObject.activeCards.Adrenaline;
    configForm['config-numSedatives'].value = configObject.activeCards.Sedatives;
    configForm['config-numCat'].value = configObject.activeCards.Cat;
    configForm['config-numMutation'].value = configObject.activeCards.Mutation;
}

function setRoleConfig(configObject) {
    let configForm = document.getElementById("lobby-gameConfig")
    configForm['config-numCaptain'].value = configObject.activeRoles.Captain;
    configForm['config-numPilot'].value = configObject.activeRoles.Pilot;
    configForm['config-numCopilot'].value = configObject.activeRoles.Copilot;
    configForm['config-numSoldier'].value = configObject.activeRoles.Soldier;
    configForm['config-numEngineer'].value = configObject.activeRoles.Engineer;
    configForm['config-numPsychologist'].value = configObject.activeRoles.Psychologist;
    configForm['config-numEO'].value = configObject.activeRoles['Executive Officer'];
    configForm['config-numMedic'].value = configObject.activeRoles.Medic;

    configForm['config-numFast'].value = configObject.activeRoles.Fast;
    configForm['config-numSurge'].value = configObject.activeRoles.Surge;
    configForm['config-numBlink'].value = configObject.activeRoles.Blink;
    configForm['config-numSilent'].value = configObject.activeRoles.Silent;
    configForm['config-numBrute'].value = configObject.activeRoles.Brute;
    configForm['config-numInvisible'].value = configObject.activeRoles.Invisible;
    configForm['config-numLurking'].value = configObject.activeRoles.Lurking;
    configForm['config-numPsychic'].value = configObject.activeRoles.Psychic;

    configForm['config-numCaptainRequired'].value = configObject.requiredRoles.Captain;
    configForm['config-numPilotRequired'].value = configObject.requiredRoles.Pilot;
    configForm['config-numCopilotRequired'].value = configObject.requiredRoles.Copilot;
    configForm['config-numSoldierRequired'].value = configObject.requiredRoles.Soldier;
    configForm['config-numEngineerRequired'].value = configObject.requiredRoles.Engineer;
    configForm['config-numPsychologistRequired'].value = configObject.requiredRoles.Psychologist;
    configForm['config-numEORequired'].value = configObject.requiredRoles['Executive Officer'];
    configForm['config-numMedicRequired'].value = configObject.requiredRoles.Medic;

    configForm['config-numFastRequired'].value = configObject.requiredRoles.Fast;
    configForm['config-numSurgeRequired'].value = configObject.requiredRoles.Surge;
    configForm['config-numBlinkRequired'].value = configObject.requiredRoles.Blink;
    configForm['config-numSilentRequired'].value = configObject.requiredRoles.Silent;
    configForm['config-numBruteRequired'].value = configObject.requiredRoles.Brute;
    configForm['config-numInvisibleRequired'].value = configObject.requiredRoles.Invisible;
    configForm['config-numLurkingRequired'].value = configObject.requiredRoles.Lurking;
    configForm['config-numPsychicRequired'].value = configObject.requiredRoles.Psychic;
}

function getGameConfig() {
    let configForm = document.getElementById("lobby-gameConfig")
    let config = {};

    config.numHumans = getConfigValue("config-numHumans")
    config.numAliens = getConfigValue("config-numAliens")

    config.numWorkingPods = getConfigValue('config-numWorkingPods')
    config.numBrokenPods = getConfigValue('config-numBrokenPods')

    config.numTurns = getConfigValue('config-numTurns')
    config.aliensRespawn = configForm['config-aliensRespawn']?.checked;

    config.activeCards = {
        "Red Card": getConfigValue('config-numRedCards'),
        "Green Card": getConfigValue('config-numGreenCards'),
        "White Card": getConfigValue('config-numWhiteCards'),
        "Teleport": getConfigValue('config-numTeleport'),
        "Clone": getConfigValue('config-numClone'),
        "Defense": getConfigValue('config-numDefense'),
        "Spotlight": getConfigValue('config-numSpotlight'),
        "Attack": getConfigValue('config-numAttack'),
        "Sensor": getConfigValue('config-numSensor'),
        "Adrenaline": getConfigValue('config-numAdrenaline'),
        "Sedatives": getConfigValue('config-numSedatives'),
        "Cat": getConfigValue('config-numCat'),
        "Mutation": getConfigValue('config-numMutation'),
    }

    config.activeStatusEffects = {
        "Armored": 2,
        "Cloned": 1
    }

    config.activeRoles = {
        "Captain": getConfigValue('config-numCaptain'),
        "Pilot": getConfigValue('config-numPilot'),
        "Copilot": getConfigValue('config-numCopilot'),
        "Soldier": getConfigValue('config-numSoldier'),
        "Psychologist": getConfigValue('config-numPsychologist'),
        "Executive Officer": getConfigValue('config-numEO'),
        "Medic": getConfigValue('config-numMedic'),
        "Engineer": getConfigValue('config-numEngineer'),

        "Fast": getConfigValue('config-numFast'),
        "Surge": getConfigValue('config-numSurge'),
        "Blink": getConfigValue('config-numBlink'),
        "Silent": getConfigValue('config-numSilent'),
        "Brute": getConfigValue('config-numBrute'),
        "Invisible": getConfigValue('config-numInvisible'),
        "Lurking": getConfigValue('config-numLurking'),
        "Psychic": getConfigValue('config-numPsychic'),
    }

    config.requiredRoles = {
        "Captain": getConfigValue('config-numCaptainRequired'),
        "Pilot": getConfigValue('config-numPilotRequired'),
        "Copilot": getConfigValue('config-numCopilotRequired'),
        "Soldier": getConfigValue('config-numSoldierRequired'),
        "Psychologist": getConfigValue('config-numPsychologistRequired'),
        "Executive Officer": getConfigValue('config-numEORequired'),
        "Medic": getConfigValue('config-numMedicRequired'),
        "Engineer": getConfigValue('config-numEngineerRequired'),

        "Fast": getConfigValue('config-numFastRequired'),
        "Surge": getConfigValue('config-numSurgeRequired'),
        "Blink": getConfigValue('config-numBlinkRequired'),
        "Silent": getConfigValue('config-numSilentRequired'),
        "Brute": getConfigValue('config-numBruteRequired'),
        "Invisible": getConfigValue('config-numInvisibleRequired'),
        "Lurking": getConfigValue('config-numLurkingRequired'),
        "Psychic": getConfigValue('config-numPsychicRequired'),
    }

    function getConfigValue(inputKey) {
        return configForm[inputKey]?.value ? parseInt(configForm[inputKey].value) : 0;
    }

    return config;
}

function updatePossible(inputName) {
    let configForm = document.getElementById("lobby-gameConfig")
    let possible = configForm[`config-${inputName}`]
    let required = configForm[`config-${inputName}Required`]

    possible.min = required.value ? parseInt(required.value) : 0
    possible.value = Math.max(possible.value ? parseInt(possible.value) : 0, possible.min ? parseInt(possible.min) : 0)
}

function checkPossible(inputName) {
    let configForm = document.getElementById("lobby-gameConfig")
    let possible = configForm[`config-${inputName}`]
    possible.value = Math.max(possible.value ? parseInt(possible.value) : 0, possible.min ? parseInt(possible.min) : 0)
}