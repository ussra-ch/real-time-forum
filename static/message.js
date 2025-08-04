import { webSocket } from "./websocket.js"
export function mesaageDiv(user, userId, receiverId) {
    // console.log();

    const body = document.querySelector('body')
    if (document.getElementById('message')) {
        document.getElementById('message').remove()
    }
    const div = document.createElement('div')
    div.id = 'message'
    div.innerHTML = `
        <div class="head">
            <h3>${user}</h3>
        </div>
        <div class="body" id="chat-body"></div>
        <form class="input-area">
            <input type="text" placeholder="..." class="chat-input" required>
            <button type="submit" class="send-btn"><i class="fa-solid fa-paper-plane"></i></button>
        </form>
    `
    const deleteButton = document.createElement('button')
    const icon = document.createElement('i');
    icon.className = 'fas fa-trash';
    deleteButton.id = 'closeConversation'
    deleteButton.appendChild(icon);
    deleteButton.appendChild(document.createTextNode(' Delete'));
    let divHeader = div.getElementsByClassName('head')[0]
    divHeader.append(deleteButton)
    body.append(div)
    //throtlle khas tkoun hna
    fetchMessages(userId, receiverId)

    div.querySelector('.input-area').addEventListener('submit', (e) => {
        e.preventDefault()
        // if (userStatus) {
        const input = div.querySelector('.chat-input')
        const message = input.value.trim()
        let chatBody = document.getElementById('chat-body')
        if (message !== "") {
            // console.log(1);
            let newMsg = document.createElement('div')
            newMsg.className = 'messageSent'
            newMsg.innerHTML = `<h3>${message}</h3>
                        <h7>${Date.now()}}</h7>`
            chatBody.append(newMsg)
            
            webSocket(userId, receiverId, input.value)
            input.value = ''
        }
        // } else {
        //     const notif = document.getElementById('not')
        //     const newNotif = document.createElement('div')
        //     newNotif.id = 'notification'
        //     // newNotif.innerHTML = `<h5> ${}`
        // }
    })

    deleteButton.addEventListener('click', () => {
        div.remove()
    })
}


function fetchMessages(userId, receiverId) {
    fetch("/api/fetchMessages")
        .then(response => response.json())
        .then(messages => {
            // console.log('dkhal lhnaa');
            // console.log(messages);
            if (!messages) {
                return
            }
            messages.reverse().forEach(message => {
                let body = document.getElementById('chat-body')
                // console.log(message);
                // console.log(message.receiverId, userId);
                // console.log(message.userId, receiverId);
                if (message.userId == userId && message.sender_id == receiverId) {
                    let newMsg = document.createElement('div')
                    newMsg.className = 'messageSent'
                    newMsg.style.background = 'blue'
                    newMsg.innerHTML = `<h3>${message.content}</h3>
            <h7>${message.time}</h7>`
                    body.append(newMsg)
                } else if (message.userId == receiverId && message.sender_id == userId) {
                    let newMsg = document.createElement('div')
                    newMsg.className = 'messageReceived'
                    newMsg.innerHTML = `<h3>${message.content}</h3>
            <h7>${message.time}</h7>`
                    body.append(newMsg)
                }
            });
        })
}
