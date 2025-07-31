export function webSockert(senderId, receiverId, messageContent) {
    const ws = new WebSocket("ws://localhost:8080/chat")


    if (!senderId || !receiverId || !messageContent) {
        return
    }

    const payload = {
        senderId,
        receiverId,
        messageContent,
    };
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
