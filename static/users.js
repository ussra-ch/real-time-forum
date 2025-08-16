import { mesaageDiv } from "./message.js";
import { ws } from "./websocket.js"
import { main } from "./main.js";
import { isAuthenticated } from "./login.js";
import { triggerUserLogout } from "./logout.js";

export function fetchUser() {
  const usern = document.getElementById('users');
  fetch('/user').then(r => r.json()).then(users => {
    usern.textContent = '';

    let on = new Set()
    if (users.onlineUsers) {
      users.onlineUsers.sort((a, b) => {
        return a.nickname.localeCompare(b.nickname);
      });

      users.onlineUsers.sort((a, b) => {
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
      const conversationButton = document.createElement('button')
      conversationButton.id = "conversationButton"
      conversationButton.innerHTML = `
      <i class="fa-solid fa-message"></i>`

      conversationButton.style.marginRight = '0'
      const div = document.createElement('div');
      div.innerHTML = `${profil} ${user.nickname}`;
      if (user.status == 'online') {
        div.style.color = 'rgb(89, 230, 187)';
      }
      div.style.display = 'flex';
      div.style.justifyContent = 'space-between';
      div.style.alignItems = 'center'
      div.style.border = '1px solid #ccc';
      div.style.padding = '8px';
      div.style.borderRadius = '70px';
      div.style.width = '60%';
      div.style.margin = '10px'
      div.style.maxWidth = '200px'
      div.style.minWidth = '120px'
      div.style.background = 'rgba(26, 35, 50, 0.95)';
      div.append(conversationButton)

      usern.appendChild(div);


      conversationButton.addEventListener('click', () => {
        isAuthenticated().then(auth => {
          if (!auth) {
            triggerUserLogout()
            main()
          } else {
            if (document.getElementById('message')) document.getElementById('message').remove()
            let isConversationOpen = { "senderId": users.UserId, "receiverId": user.userId, "isOpen": true, "type": "OpenConversation" }
            const jsonIsConversationOpen = JSON.stringify(isConversationOpen);
            ws.send(jsonIsConversationOpen);
            mesaageDiv(user.nickname, users.UserId, user.userId)
          }
        })
      })
    }
  });
}