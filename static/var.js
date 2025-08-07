export const loginDiv = document.createElement('div');
export const content = document.getElementById('content')
export let clientStatus = false
export   const ws = new WebSocket("ws://localhost:8080/chat")
export let connectedUsers = new Map();
export let notifications = 0
// export let isConversationOpen = new {}