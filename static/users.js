import { mesaageDiv } from "./message.js";
import { ws } from "./websocket.js"
import { notifications } from "./var.js";


export function fetchUser() {
  // console.log(1);

  const usern = document.getElementById('users');

  usern.innerHTML = ``
  fetch('/user').then(r => r.json()).then(users => {

    let on = new Set()
    if (users.onlineUsers) {
      users.onlineUsers.sort((a, b) => {
        return a.nickname.localeCompare(b.nickname);
      });

      users.onlineUsers.sort((a, b) => {
        console.log(a);

        return new Date(b.time) - new Date(a.time);
      });

      on = [...new Set(users.onlineUsers)];
    }

    for (const user of on) {


      if (users.UserId == user.userId) {
        continue
      }
      let profil = `<i class="fa-solid fa-user"></i>`
      if (user.photo.Valid) {
        profil = `<img src="${user.photo.String}" class="profil" class="profil" alt="Profile Picture">`
      }
      const notif = document.createElement('div')
      notif.innerHTML = `
                    <div class="notification-circle">
                    🔔
                    <div class="notification-badge" , id ="notification-circle">${notifications}</div>
                </div>
      `
      const conversationButton = document.createElement('button')
      conversationButton.id = "conversationButton"
      if (notifications == 0) {
        conversationButton.innerHTML = `
        <i class="fa-solid fa-message"></i>`
        notif.style.display = "none"
      }
      conversationButton.style.marginRight = '0'
      const div = document.createElement('div');
      div.innerHTML = `${profil} ${user.nickname}`;
      // for (const value of connectedUsers.values()){
      if (user.status == 'online') {
        // console.log(12121212);
        div.style.color = 'rgb(89, 230, 187)';
      }
      // }
      div.style.display = 'flex';
      div.style.justifyContent = 'space-between';
      div.style.alignItems = 'center'
      div.style.border = '1px solid #ccc';
      div.style.padding = '8px';
      div.style.borderRadius = '70px';
      div.style.width = '60%';
      div.style.margin = '10px'
      div.style.maxWidth = '200px'
      div.style.background = 'rgba(26, 35, 50, 0.95)';
      div.append(conversationButton)
      conversationButton.append(notif)
      usern.appendChild(div);


      conversationButton.addEventListener('click', () => {
        if (document.getElementById('message')) document.getElementById('message').remove()
        let isConversationOpen = { "senderId": users.UserId, "receiverId": user.userId, "isOpen": true, "type": "OpenConversation" }
        const jsonIsConversationOpen = JSON.stringify(isConversationOpen);
        ws.send(jsonIsConversationOpen);
        mesaageDiv(user.nickname, users.UserId, user.userId)
      })
    }
  });
}