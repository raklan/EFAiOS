typingQueue = {};
let notificationQueue = [];
let notificationCloseTimeout = null;

function typeWord(element, word) {
    if(!typingQueue[element.id]){
        typingQueue[element.id] = []
    }else{
        typingQueue[element.id].forEach(timeoutId => clearTimeout(timeoutId))
    }
    let finalInterval = 0;
    for (let i = 0; i < word.length; i++) {        
        typingQueue[element.id].push(setTimeout(typeLetter, (35 * i) + 35, element, word.charAt(i)))
        if (i == word.length - 1) {
            finalInterval = 35 + (35 * i)
        }
    }

    return finalInterval + 35;
}

function typeLetter(element, letter) {
    element.innerHTML += letter
}

function numberToLetter(num) {
    let result = '';

    while (num >= 0) {
        result = String.fromCharCode(65 + (num % 26)) + result;
        num = Math.floor(num / 26) - 1;

        if (num < 0) break;
    }

    return result;
}

function showNotification(notificationContent, notificationType) {
    var popup = document.getElementById("notification-popup");
    var content = document.getElementById("notification-content");
    var title = document.getElementById("notification-title");

    
    notificationQueue.push({
        type: notificationType,
        content: notificationContent
    })
    
    if(!notificationCloseTimeout){
        title.innerHTML = '';
        content.innerHTML = '';
        let toType = notificationQueue.shift();
        typeWord(title, toType.type);
        typeWord(content, toType.content);
    
        popup.classList.add('notification-displayed');
        notificationCloseTimeout = setTimeout(hideNotification, Math.log(Math.pow(toType.content.length, 5)) * 1000 / 3)
    }
}

function showGameOver(){
    var gameOverMsg = document.getElementById("gameover-notification");
    gameOverMsg.style.display = '';
    typeWord(document.getElementById('gameover-notification-content'), "The Game has ended!")
}

function hideNotification() {
    var popup = document.getElementById("notification-popup");    
    clearTimeout(notificationCloseTimeout);
    notificationCloseTimeout = null;

    if(notificationQueue.length > 0){
        var content = document.getElementById("notification-content");
        var title = document.getElementById("notification-title");
        title.innerHTML = '';
        content.innerHTML = '';
        let toType = notificationQueue.shift();
        typeWord(title, toType.type);
        typeWord(content, toType.content);        
        notificationCloseTimeout = setTimeout(hideNotification, Math.log(Math.pow(toType.content.length, 5)) * 1000 / 3)
    }else{
        popup.classList.remove('notification-displayed');
    }
}

function hideConfig() {
    var configpopup = document.querySelector(".config-popup-displayed")
    configpopup.classList.remove("config-popup-displayed")
    configpopup.classList.add("config-popup-hiding")
    if (!configpopup.onanimationend) {
        configpopup.onanimationend = () => {
            configpopup.classList.remove("config-popup-hiding")
        }
    }
}

function showConfig() {
    document.querySelector(".config-popup").classList.add("config-popup-displayed")
}

function configTabSwitch(newTab) {
    let generalConfigId = "config-general"
    let cardConfigId = "config-cards"
    let roleConfigId = "config-roles"
    let presetsConfigId = "config-presets"

    let generalConfig = document.getElementById(generalConfigId)
    let cardConfig = document.getElementById(cardConfigId)
    let roleConfig = document.getElementById(roleConfigId)
    let presetsConfig = document.getElementById(presetsConfigId);

    switch (newTab) {
        case generalConfigId:
            generalConfig.style.display = ''
            cardConfig.style.display = 'none';
            roleConfig.style.display = 'none';
            presetsConfig.style.display = 'none';
            break;
        case cardConfigId:
            generalConfig.style.display = 'none'
            cardConfig.style.display = '';
            roleConfig.style.display = 'none';
            presetsConfig.style.display = 'none';
            break;
        case roleConfigId:
            generalConfig.style.display = 'none'
            cardConfig.style.display = 'none';
            roleConfig.style.display = '';
            presetsConfig.style.display = 'none';
            break;
        case presetsConfigId:
            generalConfig.style.display = 'none'
            cardConfig.style.display = 'none';
            roleConfig.style.display = 'none';
            presetsConfig.style.display = 'grid';
            break;
    }
}