// import { ws } from "./var.js";
// import "./var.js";
import { fetchUser } from "./users.js";
// import { connectedUsers } from "./var.js";

// export function webSocket(senderId, receiverId, messageContent, type) {
//     if (!senderId || !receiverId || !messageContent) {
//         return
//     }

//     const payload = {
//         senderId,
//         receiverId,
//         messageContent,
//         type,
//     };
//     console.log("type is :", payload.type);

//     ws.send(JSON.stringify(payload));

// }
export var ws = null

export function initWebSocket(onMessageCallback) {
    ws = new WebSocket("ws://localhost:8080/chat")
    ws.onopen = (event) => {
        console.log("WebSocket connected");
        // const data = JSON.parse(event.data);
        console.log("data in onopen is :", event.data);
        // let notifs = document.getElementById('notification-circle')
        // notifs.textContent = data.unreadCount
        // ws.send(JSON.stringify({ type: "identify", userId: senderId }));
    };

    ws.onmessage = (event) => {
        console.log("Received:",);
        if (event.data) {
            const data = JSON.parse(event.data);
            console.log("data is :", data);
            const notifications = JSON.parse(event.data);

            if (data.type === "unreadMessage") {
                let notifs = document.getElementById('notification-circle')
                notifs.textContent = data.unreadCount
            }

            if (data.type === "message") {
                console.log("type messages :", data);
                let notifs = document.getElementById('notification-circle')
                notifs.textContent = data.Notifications
                onMessageCallback(data.messageContent);


            } else if (data.type == 'notification') {
                console.log("type notifs :", notifications);
                let notifs = document.getElementById('notification-circle')
                notifs.textContent = data.unreadCount


            }
            //  else {
            //     console.log("dkhal l else wsaaaaafi");
            //     // console.log(data.userId);
            //     console.log("data fl esle hia :", data);
            //     let notifs = document.getElementById('notification-circle')
            //     notifs.textContent = data.unreadCount
            //     fetchUser(data.userId)
            // }
        }
    };

    ws.onerror = (err) => {
        console.log('traat error : ', err);
    };

    ws.onclose = (event) => {
        console.log("WebSocket closed");
        console.log('Reason:', event.reason);
        ws.send('logout')

    };
}
