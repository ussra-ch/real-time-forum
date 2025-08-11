import { fetchUser } from "./users.js";
export var ws = null

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
            }else if (data.type === "message") {
                let notifs = document.getElementById('notification-circle')
                notifs.textContent = data.Notifications
                onMessageCallback(data.messageContent);
                fetchUser()
            } else if (data.type == 'notification' || data.type === "unreadMessage") {
                let notifs = document.getElementById('notification-circle')
                notifs.textContent = data.unreadCount
            }
        }
    };

    ws.onerror = (err) => {
        console.log('websocket error : ', err);
    };

    ws.onclose = (event) => {
        console.log("WebSocket closed");
        console.log('Reason:', event.reason);
        ws.send('logout')
    };
}
