import { mesaageDiv } from "./message.js";
export function fetchUser(id = 0) {
  const usern = document.getElementById('users');

  usern.innerHTML = ``
  fetch('/user').then(r => r.json()).then(users => {
    let on = new Set(), off = new Set()
    if (users.onlineUsers) {
      users.onlineUsers.sort((a, b) => { a.nickname - b.nickname })
      on = [...new Set(users.onlineUsers)];
    }
    // console.log(users);

    if (users.offlineUsers) {
      users.offlineUsers.sort((a, b) => { a.nickname - b.nickname }); // matnsaaaaawch diro to lowerCase
      off = [...new Set(users.offlineUsers)]
    }
    on.forEach(user => {
      console.log(user);
      let profil = `<i class="fa-solid fa-user"></i>`
      console.log(user.photo);

      if (user.photo.Valid) {
        profil = `<img src="${user.photo.String}" class="profil" class="profil" alt="Profile Picture">`
      }
      const btn = document.createElement('button')
      btn.innerHTML = `
      <i class="fa-solid fa-message"></i>`
      btn.style.marginRight = '0'
      const div = document.createElement('div');
      div.innerHTML = `${profil} ${user.nickname}`;
      console.log(id);
      if (id == user.userId) {
        
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
      div.style.background = 'rgba(26, 35, 50, 0.95)';
      div.append(btn)
      usern.appendChild(div);
      btn.addEventListener('click', () => {

        if (document.getElementById('message')) document.getElementById('message').remove()
        // console.log(user.nickname, users.UserId, user.userId);
        mesaageDiv(user.nickname, users.UserId, user.userId)
      })
    });


    off.forEach(user => {
      const btn = document.createElement('button')
      btn.innerHTML = `
      <i class="fa-solid fa-message"></i>`
      btn.style.marginRight = '0'
      const div = document.createElement('div');
      div.innerHTML = `<i class="fa-solid fa-user"></i> ${user.nickname}`;
      div.style.display = 'flex';
      div.style.justifyContent = 'space-between';
      div.style.alignItems = 'center'
      div.style.border = '1px solid #ccc';
      div.style.padding = '8px';
      div.style.borderRadius = '70px';
      div.style.width = '60%';
      div.style.maxWidth = '200px'
      div.style.margin = '10px'
      div.style.background = 'rgba(26, 35, 50, 0.95)';
      div.append(btn)
      usern.appendChild(div);
      btn.addEventListener('click', () => {
        mesaageDiv(user.nickname, users.UserId, user.userId)
      })
    });
  });
}