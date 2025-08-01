import { webSocket } from "./websocket.js"
export function mesaageDiv(user, userId, receiverId) {
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
    body.append(div)
    fetchMessages(userId, receiverId)
    div.querySelector('.input-area').addEventListener('submit', (e) => {
        e.preventDefault()
        const input = div.querySelector('.chat-input')
        const message = input.value.trim()
        let chatBody = document.getElementById('chat-body')
        if (message !== "") {
             let newMsg = document.createElement('div')
                    newMsg.className = 'messageSent'
                    newMsg.innerHTML = `<h3>${message}</h3>
                    <h7>${Date.now()}}</h7>`
                    chatBody.append(newMsg)
            webSocket(userId, receiverId, input.value)
            input.value = ''
        }
    })
}


function fetchMessages(userId, receiverId) {
    fetch("/api/fetchMessages")
        .then(response => response.json())
        .then(messages => {
            // console.log('dkhal lhnaa');
            // console.log(messages);
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
