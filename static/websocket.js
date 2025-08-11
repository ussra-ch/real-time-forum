import { fetchUser } from "./users.js";
export var ws = null

export function initWebSocket(onMessageCallback) {
    ws = new WebSocket("ws://localhost:8080/chat")
    ws.onopen = (event) => {

    };

    ws.onmessage = (event) => {


        if (event.data) {
            const data = JSON.parse(event.data);

            const notifications = JSON.parse(event.data);

            if (data.type === "unreadMessage") {

                let notifs = document.getElementById('notification-circle')
                notifs.textContent = data.unreadCount
            }
            console.log(data);

            if (data.type === 'online') {
                console.log(122);

                fetchUser()
            } else if (data.type === 'offline') {
                console.log(122);

                fetchUser()
            }
            if (data.type === "message") {

                let notifs = document.getElementById('notification-circle')
                notifs.textContent = data.Notifications
                onMessageCallback(data.messageContent);
                fetchUser()


            } else if (data.type == 'notification') {


                let notifs = document.getElementById('notification-circle')
                notifs.textContent = data.unreadCount


            }
        }
    };

    ws.onerror = (err) => {

    };

    ws.onclose = (event) => {

    };
}
