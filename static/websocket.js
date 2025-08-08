import { ws } from "./var.js";
import { fetchUser } from "./users.js";
// import { connectedUsers } from "./var.js";

export function webSocket(senderId, receiverId, messageContent, seen, isOpen, type) {
    if (!senderId || !receiverId || !messageContent) {
        return
    }

    const payload = {
        senderId,
        receiverId,
        messageContent,
        seen,
        isOpen,
        type,
    };

    ws.send(JSON.stringify(payload));

}

export function initWebSocket(onMessageCallback) {
    ws.onopen = () => {
        console.log("WebSocket connected");
        // ws.send(JSON.stringify({ type: "identify", userId: senderId }));
    };

    ws.onmessage = (event) => {
        console.log("Received:",);
        const data = JSON.parse(event.data);
        if (data.type === "message") {
            const notifications = JSON.parse(event.data);
            console.log("type messages, w notifs huma :", notifications);
            let notifs = document.getElementById('notification-circle')
            notifs.textContent = notifications.unreadCount
            onMessageCallback(data.content);
        } else if (data.type == 'notification') {
            const notifications = JSON.parse(event.data);
            console.log("type notifs :", notifications);
            let notifs = document.getElementById('notification-circle')
            notifs.textContent = notifications.unreadCount
            // console.log(notifs);
        } else {
            console.log("dkhal l else wsaaaaafi");
            console.log(data.userId);
            
            fetchUser(data.userId)
        }
    };

    ws.onerror = (err) => {
        console.log('traat error : ', err);
    };

    ws.onclose = (event) => {
        console.log("WebSocket closed");
        console.log('Reason:', event.reason);
        ws.send('logut')
    
    };
}
