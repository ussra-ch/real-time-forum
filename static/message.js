import { ws } from "./websocket.js"
import { fetchUser } from "./users.js"

function throttle(func, delay) {
    let lastCall = 0;
    return function (...args) {
        const now = Date.now();
        if (now - lastCall >= delay) {
            lastCall = now;
            func.apply(this, args);
        }
    };
}
export function mesaageDiv(user, userId, receiverId) {
    const body = document.querySelector('body')
    if (document.getElementById('message')) {
        document.getElementById('message').remove()
    }
    var done = true
    const conversation = document.createElement('div')
    conversation.id = 'message'
    conversation.innerHTML = `
        <div class="head">
            <h3>${user}</h3>
        </div>
        <div class="body" id="chat-body"></div>
        <div id="footer"></div>
        <form class="input-area">
            <input type="text" placeholder="..." class="chat-input" required>
            <button type="submit" class="send-btn"><i class="fa-solid fa-paper-plane"></i></button>
        </form>
    `

    /////////////////////// Delete conversation Button
    // const deleteButton = document.createElement('button')
    // const icon = document.createElement('i');
    // icon.className = 'fas fa-trash';
    // deleteButton.id = 'closeConversation'
    // deleteButton.appendChild(icon);
    // deleteButton.appendChild(document.createTextNode(' Delete'));
    // let divHeader = conversation.getElementsByClassName('head')[0]
    // divHeader.append(deleteButton)
    body.append(conversation)
    ///////////////////////


    //throtlle khas tkoun hna
    const container = document.getElementById('chat-body')
    let offset = 0;
    const limit = 10;
    fetchMessages(userId, receiverId, offset, limit, user)
    let lastCall = 0;
    const delay = 500;

    container.addEventListener("scroll", () => {
        if (container.scrollTop === 0) {
            const now = Date.now();
            const canCall = now - lastCall >= delay;

            if (canCall) {
                lastCall = now;
                offset += limit;
                throttle(fetchMessages, 500)(userId, receiverId, offset, limit, user);
                //fetchMessages(userId, receiverId, offset, limit, user);
            }
        }
    });



    conversation.querySelector('.input-area').addEventListener('submit', (e) => {
        e.preventDefault()
        // webSocket(userId, receiverId, "", true, true, "conversation")
        const input = conversation.querySelector('.chat-input')
        const message = input.value.trim()
        let chatBody = document.getElementById('chat-body')

        if (message !== "") {
            // console.log(1);
            let newMsg = document.createElement('div')
            newMsg.className = 'messageReceived'
            let msgContent = document.createElement('h3')
            msgContent.textContent = message
            let timestamp = document.createElement('h7')
            timestamp.textContent = formatDate(Date.now())
            newMsg.appendChild(msgContent)
            newMsg.appendChild(timestamp)
            chatBody.append(newMsg)

            const payload = {
                "senderId": userId,
                "receiverId": receiverId,
                "messageContent": input.value,
                "type": "message",
            };
            // console.log("type is :", payload.type);

            ws.send(JSON.stringify(payload));
            fetchUser()
            // webSocket(userId, receiverId, input.value, "message")
            input.value = ''
            const container = document.getElementById('chat-body')
            container.scrollTop = container.scrollHeight;
        }
    })
    conversation.querySelector('.input-area').addEventListener('input', (e) => {
        const payload = {
            "senderId": userId,
            "receiverId": receiverId,
            "type": "typing",
        };

        ws.send(JSON.stringify(payload));

    })

    window.addEventListener('click', (e) => {

        if (conversation && !conversation.contains(e.target) && !done) {
            console.log(1);
            let isConversationOpen = {
                senderId: userId,
                receiverId: receiverId,
                type: "CloseConversation"
            }
            ws.send(JSON.stringify(isConversationOpen));
            conversation.remove();
        }
        done=false
    });

}

function fetchMessages(userId, receiverId, offset, limit, name) {
    const body = document.getElementById('chat-body');
    var previousScrollHeight = body.scrollHeight;

    fetch(`/api/fetchMessages?offset=${offset}&limit=${limit}&sender=${receiverId} `, {
        method: 'GET'
    })
        .then(response => response.json())
        .then(messages => {

            if (!messages) {
                return
            }


            for (const message of messages) {
                if (message.content != "") {
                    if (message.userId == userId && message.sender_id == receiverId) {
                        let newMsg = document.createElement('div')
                        newMsg.className = 'messageSent'

                        newMsg.innerHTML = `
                                            <div class="messagProfil">
                                                <div class="profile">
                                                </div>
                                                   <h7>${name}</h7>
                                            </div>
                                            <h3>${message.content}</h3>
                                            <h7>${formatDate(message.time)}</h7>`
                        body.prepend(newMsg)
                        if (document.querySelector('.profile')) {

                            if (message.photo) {
                                document.querySelector('.profile').style.backgroundImage = `url(${message.photo})`;
                            } else {
                                document.querySelector('.profile').innerHTML = `
                                    <i class="fa-solid fa-user"></i>
                                `
                            }
                        }
                    } else if (message.userId == receiverId && message.sender_id == userId) {
                        let newMsg = document.createElement('div')
                        newMsg.className = 'messageReceived'
                        newMsg.innerHTML = `
                                            <h3>${message.content}</h3>
                                            <h7>${formatDate(message.time)}</h7>`
                        body.prepend(newMsg)
                    }
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