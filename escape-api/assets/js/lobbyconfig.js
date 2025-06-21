const GAME_CONFIG_DEFAULT = {
    workingPods: 4,
    brokenPods: 1,

    numTurns: 40,
    aliensRespawn: false,

    red: 24,
    green: 26,
    silent: 4,
    adrenaline: 3,
    attack: 1,
    cat: 2,
    clone: 1,
    defense: 1,
    mutation: 1,
    sedatives: 1,
    sensor: 1,
    spotlight: 2,
    teleport: 1,

    numCaptain: 1,
    numPilot: 1,
    numCopilot: 1,
    numSoldier: 1,
    numPsychologist: 1,
    numEO: 1,
    numMedic: 1,
    numEngineer: 1,

    numFast: 1,
    numSurge: 1,
    numBlink: 1,
    numSilent: 1,
    numBrute: 1,
    numInvisible: 1,
    numLurking: 1,
    numPsychic: 1
}

const GAME_CONFIG_BASIC = {
    workingPods: 4,
    brokenPods: 1,

    numTurns: 40,
    aliensRespawn: false,

    red: 24,
    green: 26,
    silent: 18,
    adrenaline: 0,
    attack: 0,
    cat: 0,
    clone: 0,
    defense: 0,
    mutation: 0,
    sedatives: 0,
    sensor: 0,
    spotlight: 0,
    teleport: 0,

    numCaptain: 1,
    numPilot: 1,
    numCopilot: 1,
    numSoldier: 1,
    numPsychologist: 1,
    numEO: 1,
    numMedic: 1,
    numEngineer: 1,

    numFast: 1,
    numSurge: 1,
    numBlink: 1,
    numSilent: 1,
    numBrute: 1,
    numInvisible: 1,
    numLurking: 1,
    numPsychic: 1
}

const GAME_CONFIG_HIDEANDSEEK = {
    workingPods: 4,
    brokenPods: 1,

    numTurns: 40,
    aliensRespawn: false,

    red: 24,
    green: 0,
    silent: 4,
    adrenaline: 0,
    attack: 1,
    cat: 2,
    clone: 1,
    defense: 1,
    mutation: 1,
    sedatives: 1,
    sensor: 1,
    spotlight: 2,
    teleport: 1,

    numCaptain: 1,
    numPilot: 1,
    numCopilot: 1,
    numSoldier: 1,
    numPsychologist: 1,
    numEO: 1,
    numMedic: 1,
    numEngineer: 1,

    numFast: 1,
    numSurge: 1,
    numBlink: 1,
    numSilent: 1,
    numBrute: 1,
    numInvisible: 1,
    numLurking: 1,
    numPsychic: 1
}

const GAME_PRESETS = {
    Tabletop: GAME_CONFIG_DEFAULT,
    Basic: GAME_CONFIG_BASIC
}

function setAllConfigAsDefault() {
    setGeneralConfigAsDefault();
    setCardConfigAsDefault();
    setRoleConfigAsDefault();
}

function setGeneralConfigAsDefault() {
    setGeneralConfig(GAME_CONFIG_DEFAULT)
}

function setCardConfigAsDefault() {
    setCardConfig(GAME_CONFIG_DEFAULT)
}

function setRoleConfigAsDefault() {
    setRoleConfig(GAME_CONFIG_DEFAULT)
}

function setConfigForm(configObject) {
    setGeneralConfig(configObject);
    setCardConfig(configObject);
    setRoleConfig(configObject);
}

function setGeneralConfig(configObject) {
    let configForm = document.getElementById("lobby-gameConfig")

    configForm['config-numHumans'].value = 0;
    configForm['config-numAliens'].value = 0;

    configForm['config-numWorkingPods'].value = configObject.workingPods;
    configForm['config-numBrokenPods'].value = configObject.brokenPods;

    configForm['config-numTurns'].value = configObject.numTurns;
    configForm['config-aliensRespawn'].checked = configObject.aliensRespawn
}

function setCardConfig(configObject) {
    let configForm = document.getElementById("lobby-gameConfig")

    configForm['config-numRedCards'].value = configObject.red;
    configForm['config-numGreenCards'].value = configObject.green;
    configForm['config-numWhiteCards'].value = configObject.silent;

    configForm['config-numTeleport'].value = configObject.teleport;
    configForm['config-numClone'].value = configObject.clone;
    configForm['config-numDefense'].value = configObject.defense;

    configForm['config-numSpotlight'].value = configObject.spotlight;
    configForm['config-numAttack'].value = configObject.attack;
    configForm['config-numSensor'].value = configObject.sensor;

    configForm['config-numAdrenaline'].value = configObject.adrenaline;
    configForm['config-numSedatives'].value = configObject.sedatives;
    configForm['config-numCat'].value = configObject.cat;
    configForm['config-numMutation'].value = configObject.mutation;
}

function setRoleConfig(configObject) {
    let configForm = document.getElementById("lobby-gameConfig")
    configForm['config-numCaptain'].value = configObject.numCaptain;
    configForm['config-numPilot'].value = configObject.numPilot;
    configForm['config-numCopilot'].value = configObject.numCopilot;
    configForm['config-numSoldier'].value = configObject.numSoldier;
    configForm['config-numEngineer'].value = configObject.numEngineer;
    configForm['config-numPsychologist'].value = configObject.numPsychologist;
    configForm['config-numEO'].value = configObject.numEO;
    configForm['config-numMedic'].value = configObject.numMedic;

    configForm['config-numFast'].value = configObject.numFast;
    configForm['config-numSurge'].value = configObject.numSurge;
    configForm['config-numBlink'].value = configObject.numSilent;
    configForm['config-numSilent'].value = configObject.numSilent;
    configForm['config-numBrute'].value = configObject.numBrute;
    configForm['config-numInvisible'].value = configObject.numInvisible;
    configForm['config-numLurking'].value = configObject.numLurking;
    configForm['config-numPsychic'].value = configObject.numPsychic;
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