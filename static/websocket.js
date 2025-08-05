import { ws } from "./var.js";
import { fetchUser } from "./users.js";
// import { connectedUsers } from "./var.js";

export function webSocket(senderId, receiverId, messageContent, seen) {

    if (!senderId || !receiverId || !messageContent) {
        return
    }

    const payload = {
        senderId,
        receiverId,
        messageContent,
        seen,
        // userStatus,
    };
    // console.log(payload);


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
            // console.log(22);
            // fetchUser(data.userId)
            onMessageCallback(event.data);
        } else {
            fetchUser(data.userId)
        }
    };

    ws.onerror = (err) => {
        console.log('traat error : ', err);
    };

    ws.onclose = (event) => {
        console.log("WebSocket closed");
        console.log('Reason:', event.reason);
    };
}
