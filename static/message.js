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
    const container = document.getElementById('chat-body')
    let offset = 0;
    const limit = 10;
    fetchMessages(userId, receiverId, offset, limit)

    container.addEventListener("scroll", () => {

        if (container.scrollTop === 0) {
            console.log(offset, limit);
            offset += limit;
            fetchMessages(userId, receiverId, offset, limit)
        }
    });
    div.querySelector('.input-area').addEventListener('submit', (e) => {
        e.preventDefault()
        // if (userStatus) {
        const input = div.querySelector('.chat-input')
        const message = input.value.trim()
        let chatBody = document.getElementById('chat-body')
        if (message !== "") {
            // console.log(1);
            let newMsg = document.createElement('div')
            newMsg.className = 'messageReceived'
            newMsg.innerHTML = `<h3>${message}</h3>
                        <h7>${formatDate(Date.now())}</h7>`
            chatBody.append(newMsg)
            webSocket(userId, receiverId, input.value)
            input.value = ''
            const container = document.getElementById('chat-body')
            container.scrollTop = container.scrollHeight;
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


function fetchMessages(userId, receiverId, offset, limit) {
    const body = document.getElementById('chat-body');
    var previousScrollHeight = body.scrollHeight;
    console.log(1);


    fetch(`/api/fetchMessages?offset=${offset}&limit=${limit}&sender=${receiverId} `, {
        method: 'GET'
    })
        .then(response => response.json())
        .then(messages => {
            if (!messages) {
                return
            }
            for (const message of messages) {
                if (message.userId == userId && message.sender_id == receiverId) {
                    let newMsg = document.createElement('div')
                    newMsg.className = 'messageSent'
                    newMsg.style.background = 'blue'
                    newMsg.innerHTML = `<h3>${message.content}</h3>
                                        <h7>${formatDate(message.time)}</h7>`
                    body.prepend(newMsg)
                } else if (message.userId == receiverId && message.sender_id == userId) {
                    let newMsg = document.createElement('div')
                    newMsg.className = 'messageReceived'
                    newMsg.innerHTML = `<h3>${message.content}</h3>
                                        <h7>${formatDate(message.time)}</h7>`
                    body.prepend(newMsg)
                }
            }
            const newScrollHeight = body.scrollHeight;
            body.scrollTop += (newScrollHeight - previousScrollHeight);
        })
}
export function formatDate(timestampInSeconds) {
    const isoString = timestampInSeconds;
    const date = new Date(isoString);
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    const time = `${hours}:${minutes}`;
    return time;
}