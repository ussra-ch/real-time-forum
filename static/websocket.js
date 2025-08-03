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


    ws.send(JSON.stringify(payload));

}
export function initWebSocket(onMessageCallback) {
    ws.onopen = () => {
        console.log("WebSocket connected");
        connectedUsers.set(id, 'online');
        // ws.send(JSON.stringify({ type: "identify", userId: senderId }));
    };

    ws.onmessage = (event) => {
        console.log("Received:",);

        const data = JSON.parse(event.data);
        if (data.type === "message") {
            console.log(22);
            
            onMessageCallback(event.data);
        } else {
            fetchUser(data.userId)
        }
    };

    ws.onerror = (err) => {

    };

    ws.onclose = () => {
        console.log("WebSocket closed");
    };
}
