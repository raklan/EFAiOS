function typeWord(element, word){        
    let finalInterval = 0;
    for(let i = 0; i < word.length; i++){
        setTimeout(typeLetter, (35 * i) + 35, element, word.charAt(i))
        if(i == word.length - 1){
            finalInterval = 35 + (35 * i)
        }
    }

    return finalInterval + 35;
}

function typeLetter(element, letter){
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
