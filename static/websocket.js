import { fetchUser } from "./users.js"
export var ws = null
let lastCall = 0;
let typingTimeout;

function typingInProgress(Id) {
    const chat = document.getElementById('footer');
    if (!chat) return;
    const userID = document.getElementById('message_id')
    console.log(userID.value,Id);
    
    if (userID.value != Id) {
        return
    }
    let typingEl = document.getElementById('typing');
    if (!typingEl) {
        const div = document.createElement('div');
        div.id = 'typing';
        div.innerHTML = `
            <div id="typing-indicator">
              <span></span>
              <span></span>
             <span></span>
            </div>`;
        chat.append(div);
        typingEl = div;
    }

    chat.scrollTop = chat.scrollHeight;

    if (typingTimeout) clearTimeout(typingTimeout);

    typingTimeout = setTimeout(() => {
        const el = document.getElementById('typing');
        if (el) el.remove();
    }, 1000);
}

export function initWebSocket(onMessageCallback) {
    ws = new WebSocket("ws://localhost:8080/chat")

    ws.onopen = (event) => {
        console.log("WebSocket connected");
    };

    ws.onmessage = (event) => {
        if (event.data) {

            const data = JSON.parse(event.data);

            if (data.type === 'online' || data.type === 'offline') {
                fetchUser()
            } else if (data.type === "message") {
              
                
                let notifs = document.getElementById('notification-circle')
                notifs.textContent = data.Notifications
                onMessageCallback(data.messageContent,data.name);

            } else if (data.type == 'notification' || data.type === "unreadMessage") {
                let notifs = document.getElementById('notification-circle')
                notifs.textContent = data.unreadCount

            } else if (data.type == 'typing') {
                typingInProgress(data.sender)
            }
        }
    };

    ws.onerror = (err) => {
        console.log('websocket error : ', err);
    };

    ws.onclose = (event) => {

    };
}
