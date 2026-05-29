function setConfigFormFromObject(configObject) {
    setTeamConfig(configObject);
    setCardConfig(configObject);
    setRoleConfig(configObject);
    setModifierConfig(configObject);
}

function setConfigFormFromString(configString) {
    const configObject = JSON.parse(configString);
    setConfigFormFromObject(configObject);
}

function setModifierConfig(configObject) {
    let configForm = document.getElementById("lobby-gameConfig")       

    setConfigInputValue('config-numTurns', configObject.modifiers.numTurns);
    configForm['config-aliensRespawn'].checked = configObject.modifiers.aliensRespawn;
    configForm['config-autoTurnEnd'].checked = configObject.modifiers.autoTurnEnd;
    configForm['config-survivalMode'].checked = configObject.modifiers.survivalMode;
}

function setCardConfig(configObject) {
    setConfigInputValue('config-numWorkingPods', configObject.numWorkingPods);
    setConfigInputValue('config-numBrokenPods', configObject.numBrokenPods);

    setConfigInputValue('config-numRedCards', configObject.activeCards['Red Card']);
    setConfigInputValue('config-numGreenCards', configObject.activeCards['Green Card']);
    setConfigInputValue('config-numWhiteCards', configObject.activeCards['White Card']);

    setConfigInputValue('config-numTeleport', configObject.activeCards.Teleport);
    setConfigInputValue('config-numClone', configObject.activeCards.Clone);
    setConfigInputValue('config-numDefense', configObject.activeCards.Defense);

    setConfigInputValue('config-numSpotlight', configObject.activeCards.Spotlight);
    setConfigInputValue('config-numAttack', configObject.activeCards.Attack);
    setConfigInputValue('config-numSensor', configObject.activeCards.Sensor);

    setConfigInputValue('config-numAdrenaline', configObject.activeCards.Adrenaline);
    setConfigInputValue('config-numSedatives', configObject.activeCards.Sedatives);
    setConfigInputValue('config-numCat', configObject.activeCards.Cat);
    setConfigInputValue('config-numMutation', configObject.activeCards.Mutation);

    setConfigInputValue('config-numHidingSpot', configObject.activeCards['Hiding Spot']);
    setConfigInputValue('config-numCloakingDevice', configObject.activeCards['Cloaking Device']);
    setConfigInputValue('config-numEngineeringManual', configObject.activeCards['Engineering Manual']);
    setConfigInputValue('config-numNoisemaker', configObject.activeCards.Noisemaker);
}

function setRoleConfig(configObject) {
    let configForm = document.getElementById("lobby-gameConfig")    

    setConfigInputValue('config-numCaptain', configObject.activeRoles.Captain);
    setConfigInputValue('config-numPilot', configObject.activeRoles.Pilot);
    setConfigInputValue('config-numCopilot', configObject.activeRoles.Copilot);
    setConfigInputValue('config-numSoldier', configObject.activeRoles.Soldier);
    setConfigInputValue('config-numEngineer', configObject.activeRoles.Engineer);
    setConfigInputValue('config-numPsychologist', configObject.activeRoles.Psychologist);
    setConfigInputValue('config-numEO', configObject.activeRoles['Executive Officer']);
    setConfigInputValue('config-numMedic', configObject.activeRoles.Medic);

    setConfigInputValue('config-numFast', configObject.activeRoles.Fast);
    setConfigInputValue('config-numSurge', configObject.activeRoles.Surge);
    setConfigInputValue('config-numBlink', configObject.activeRoles.Blink);
    setConfigInputValue('config-numSilent', configObject.activeRoles.Silent);
    setConfigInputValue('config-numBrute', configObject.activeRoles.Brute);
    setConfigInputValue('config-numInvisible', configObject.activeRoles.Invisible);
    setConfigInputValue('config-numLurking', configObject.activeRoles.Lurking);
    setConfigInputValue('config-numPsychic', configObject.activeRoles.Psychic);

    setConfigInputValue('config-numCaptainRequired', configObject.requiredRoles.Captain);
    setConfigInputValue('config-numPilotRequired', configObject.requiredRoles.Pilot);
    setConfigInputValue('config-numCopilotRequired', configObject.requiredRoles.Copilot);
    setConfigInputValue('config-numSoldierRequired', configObject.requiredRoles.Soldier);
    setConfigInputValue('config-numEngineerRequired', configObject.requiredRoles.Engineer);
    setConfigInputValue('config-numPsychologistRequired', configObject.requiredRoles.Psychologist);
    setConfigInputValue('config-numEORequired', configObject.requiredRoles['Executive Officer']);
    setConfigInputValue('config-numMedicRequired', configObject.requiredRoles.Medic);

    setConfigInputValue('config-numFastRequired', configObject.requiredRoles.Fast);
    setConfigInputValue('config-numSurgeRequired', configObject.requiredRoles.Surge);
    setConfigInputValue('config-numBlinkRequired', configObject.requiredRoles.Blink);
    setConfigInputValue('config-numSilentRequired', configObject.requiredRoles.Silent);
    setConfigInputValue('config-numBruteRequired', configObject.requiredRoles.Brute);
    setConfigInputValue('config-numInvisibleRequired', configObject.requiredRoles.Invisible);
    setConfigInputValue('config-numLurkingRequired', configObject.requiredRoles.Lurking);
    setConfigInputValue('config-numPsychicRequired', configObject.requiredRoles.Psychic);

    setConfigInputValue('config-numScout', configObject.activeRoles.Scout);
    setConfigInputValue('config-numCommunicationsOfficer', configObject.activeRoles['Communications Officer']);

    setConfigInputValue('config-numTracker', configObject.activeRoles.Tracker);
    setConfigInputValue('config-numCalling', configObject.activeRoles.Calling);

    setConfigInputValue('config-numScoutRequired', configObject.requiredRoles.Scout);
    setConfigInputValue('config-numCommunicationsOfficerRequired', configObject.requiredRoles['Communications Officer']);

    setConfigInputValue('config-numTrackerRequired', configObject.requiredRoles.Tracker);
    setConfigInputValue('config-numCallingRequired', configObject.requiredRoles.Calling);
}

function setTeamConfig(configObject){
    setConfigInputValue('config-numHumans', configObject.numHumans);
    setConfigInputValue('config-numAliens', configObject.numAliens);     
}

function getGameConfig() {
    let configForm = document.getElementById("lobby-gameConfig")
    let config = {}

    config.numHumans = getConfigValue("config-numHumans")
    config.numAliens = getConfigValue("config-numAliens")

    config.teamAssignments = {}
    let teamAssignments = document.querySelectorAll('.player-team-assignment')
    for(let assignment of teamAssignments){
        if(configForm[assignment.id].value?.length > 0){
            config.teamAssignments[assignment.id.replace('team-assignment-', '')] = configForm[assignment.id].value
        }
    }

    config.numWorkingPods = getConfigValue('config-numWorkingPods')
    config.numBrokenPods = getConfigValue('config-numBrokenPods')

    config.modifiers = {
        numTurns: getConfigValue('config-numTurns'),
        aliensRespawn: configForm['config-aliensRespawn']?.checked,
        autoTurnEnd: configForm['config-autoTurnEnd']?.checked,
        survivalMode: configForm['config-survivalMode']?.checked,
    }    

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

        "Hiding Spot": getConfigValue('config-numHidingSpot'),
        "Cloaking Device": getConfigValue('config-numCloakingDevice'),
        "Engineering Manual": getConfigValue('config-numEngineeringManual'),
        "Noisemaker": getConfigValue('config-numNoisemaker')
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

        'Scout': getConfigValue('config-numScout'),
        'Communications Officer': getConfigValue('config-numCommunicationsOfficer'),

        'Tracker': getConfigValue('config-numTracker'),
        'Calling': getConfigValue('config-numCalling'),
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

        'Scout': getConfigValue('config-numScoutRequired'),
        'Communications Officer': getConfigValue('config-numCommunicationsOfficerRequired'),

        'Tracker': getConfigValue('config-numTrackerRequired'),
        'Calling': getConfigValue('config-numCallingRequired'),
    }

    function getConfigValue(inputKey) {
        return configForm[inputKey]?.value ? parseInt(configForm[inputKey].value) : 0;
    }

    return config;
}

function setConfigInputValue(inputName, val) {
    let configForm = document.getElementById("lobby-gameConfig")
    configForm[inputName].value = val ?? 0;
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