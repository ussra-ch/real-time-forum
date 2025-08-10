export const loginDiv = document.createElement('div');
export const content = document.getElementById('content')
export let clientStatus = false
// export var ws = null
export let connectedUsers = new Map();
export let notifications = 0
// export let isConversationOpen = new {}

// export function initWebSocket(onMessageCallback) {
//     ws.onopen = (event) => {
//         console.log("WebSocket connected");
//         const data = JSON.parse(event.data);
//         console.log("data in onopen is :", data);
//         let notifs = document.getElementById('notification-circle')
//         notifs.textContent = data.unreadCount
//         // ws.send(JSON.stringify({ type: "identify", userId: senderId }));
//     };

//     // ws.onmessage = (event) => {
//     //     console.log("Received:",);
//     //     if (event.data) {
//     //         const data = JSON.parse(event.data);
//     //         console.log("data is :", data);
//     //         const notifications = JSON.parse(event.data);


//     //         if (data.type === "message") {
//     //             console.log("type messages :", data);
//     //             let notifs = document.getElementById('notification-circle')
//     //             notifs.textContent = data.Notifications
//     //             onMessageCallback(data.messageContent);


//     //         } else if (data.type == 'notification') {
//     //             console.log("type notifs :", notifications);
//     //             let notifs = document.getElementById('notification-circle')
//     //             notifs.textContent = data.unreadCount


//     //         } else {
//     //             console.log("dkhal l else wsaaaaafi");
//     //             // console.log(data.userId);
//     //             console.log("data fl esle hia :", data);
//     //             let notifs = document.getElementById('notification-circle')
//     //             notifs.textContent = data.notifications
//     //             fetchUser(data.userId)
//     //         }
//     //     }
//     // };

//     ws.onerror = (err) => {
//         console.log('traat error : ', err);
//     };

//     ws.onclose = (event) => {
//         console.log("WebSocket closed");
//         console.log('Reason:', event.reason);
//         ws.send('logout')

//     };
// }

// initWebSocket()
// ws.onopen = ()=>{
//     console.log("var file on open");
// }