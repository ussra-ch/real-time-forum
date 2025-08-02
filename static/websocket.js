import { ws } from "./var.js";
export function webSocket(senderId, receiverId, messageContent) {
    console.log("js webSocket: 1111");
    
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
    
    ws.onopen = () => {
        ws.send(JSON.stringify(payload));
        console.log('Connected!')
    };

    ws.onmessage = (event) => {
        let chatBody = document.getElementById('chat-body')
        let newMsg = document.createElement('div')
        newMsg.innerHTML = `<h3>${event.data}</h3>`
        chatBody.append(newMsg)
    }
    ws.onerror = (err) => console.error('Error:', err);
    ws.onclose = () => console.log('Closed');
}
