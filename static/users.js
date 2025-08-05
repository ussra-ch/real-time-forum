import { mesaageDiv } from "./message.js";
// import { connectedUsers } from "./var.js";

export function fetchUser() {
  // console.log(connectedUsers);

  const usern = document.getElementById('users');

  usern.innerHTML = ``
  fetch('/user').then(r => r.json()).then(users => {
    console.log(users);
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
      console.log(user.time);

      if (users.UserId == user.userId) {
        continue
      }
      let profil = `<i class="fa-solid fa-user"></i>`
      if (user.photo.Valid) {
        profil = `<img src="${user.photo.String}" class="profil" class="profil" alt="Profile Picture">`
      }
      const btn = document.createElement('button')
      btn.innerHTML = `
      <i class="fa-solid fa-message"></i>`
      btn.style.marginRight = '0'
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
      div.append(btn)
      usern.appendChild(div);


      btn.addEventListener('click', () => {
        if (document.getElementById('message')) document.getElementById('message').remove()
        // console.log(user.nickname, users.id, user.userId);

        mesaageDiv(user.nickname, users.UserId, user.userId)
      })
    }


    // off.forEach(user => {
    //   const btn = document.createElement('button')
    //   btn.innerHTML = `
    //   <i class="fa-solid fa-message"></i>`
    //   btn.style.marginRight = '0'
    //   const div = document.createElement('div');
    //   div.innerHTML = `<i class="fa-solid fa-user"></i> ${user.nickname}`;
    //   div.style.display = 'flex';
    //   div.style.justifyContent = 'space-between';
    //   div.style.alignItems = 'center'
    //   div.style.border = '1px solid #ccc';
    //   div.style.padding = '8px';
    //   div.style.borderRadius = '70px';
    //   div.style.width = '60%';
    //   div.style.maxWidth = '200px'
    //   div.style.margin = '10px'
    //   div.style.background = 'rgba(26, 35, 50, 0.95)';
    //   div.append(btn)
    //   usern.appendChild(div);
    //   btn.addEventListener('click', () => {
    //     mesaageDiv(user.nickname, users.UserId, user.userId)
    //   })
    // });
  });
}