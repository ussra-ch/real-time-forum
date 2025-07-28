import { mesaageDiv } from "./messag.js";
export function fetchUser() {
  const usern = document.getElementById('users');

  usern.innerHTML = ``
  fetch('/user').then(r => r.json()).then(users => {
    let on = new Set(), off = new Set()
    if (users.onlineUsers) {
      users.onlineUsers.sort((a, b) => { a - b })
      on = [...new Set(users.onlineUsers)];

    }
    if (users.offlineUsers) {

      users.offlineUsers.sort((a, b) => { a - b });
      off = [...new Set(users.offlineUsers)]
    }
    on.forEach(user => {
      const btn = document.createElement('button')
      btn.innerHTML = `
      <i class="fa-solid fa-message"></i>`
      btn.style.marginRight = '0'
      const div = document.createElement('div');
      div.innerHTML = `<i class="fa-solid fa-user"></i> ${user}`;
      div.style.color = 'rgb(89, 230, 187)';
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
        if (document.getElementById('messag')) document.getElementById('messag').remove()
        mesaageDiv(user)
      })
    });

    off.forEach(user => {
      const btn = document.createElement('button')
      btn.innerHTML = `
      <i class="fa-solid fa-message"></i>`
      btn.style.marginRight = '0'
      const div = document.createElement('div');
      div.innerHTML = `<i class="fa-solid fa-user"></i> ${user}`;
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
        //if (document.getElementById('messag')) document.getElementById('messag').remove()
        mesaageDiv(user)
      })
    });
  });
}