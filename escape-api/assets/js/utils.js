function typeWord(element, word) {
    let finalInterval = 0;
    for (let i = 0; i < word.length; i++) {
        setTimeout(typeLetter, (35 * i) + 35, element, word.charAt(i))
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

    title.innerHTML = '';
    content.innerHTML = '';

    typeWord(title, notificationType)
    typeWord(content, notificationContent)

    popup.classList.add('notification-displayed')
}

function hideNotification() {
    var popup = document.getElementById("notification-popup");
    popup.classList.remove('notification-displayed')
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

    let generalConfig = document.getElementById(generalConfigId)
    let cardConfig = document.getElementById(cardConfigId)
    let roleConfig = document.getElementById(roleConfigId)

    switch (newTab) {
        case generalConfigId:
            generalConfig.style.display = ''
            cardConfig.style.display = 'none';
            roleConfig.style.display = 'none';
            break;
        case cardConfigId:
            generalConfig.style.display = 'none'
            cardConfig.style.display = '';
            roleConfig.style.display = 'none';
            break;
        case roleConfigId:
            generalConfig.style.display = 'none'
            cardConfig.style.display = 'none';
            roleConfig.style.display = '';
            break;
    }
}