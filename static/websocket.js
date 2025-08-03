import { ws } from "./var.js";
import { fetchUser } from "./users.js"; 
// import { connectedUsers } from "./var.js";

export function webSocket(senderId, receiverId, messageContent) {

    if (!senderId || !receiverId || !messageContent) {
        return
    }

    const payload = {
        senderId,
        receiverId,
        messageContent,
        // userStatus,
    };
    // console.log(payload);

    if (ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify(payload));
    } else {
        console.warn("WebSocket not open. Message not sent.");
    }
}
export function initWebSocket(onMessageCallback) {
    ws.onopen = () => {
        console.log("WebSocket connected");
        connectedUsers.set(id, 'online');
        // ws.send(JSON.stringify({ type: "identify", userId: senderId }));
    };

    ws.onmessage = (event) => {
        console.log("Received:", event);
        // if (typeof onMessageCallback === 'function') {
            // onMessageCallback(event.data);
        // }

        const data = JSON.parse(event.data);
        console.log('-------', event.data);        
        if (data.type === "online") {
            console.log(1212);
            
            fetchUser(data.userId)
        }
    };

    ws.onerror = (err) => {

    };

    ws.onclose = () => {
        console.log("WebSocket closed");
    };
}
