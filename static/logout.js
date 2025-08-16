import { login, logindiv } from "./login.js";
import { fetchUser } from "./users.js";
import { ws } from "./websocket.js";
// import { loginDiv } from "./var.js";
// import { connectedUsers } from "./var.js";



//katloggouti (katfetshi data mn backend bach katms7 session)
export function logout() {
  const Logout = document.getElementById('logout');

  Logout.addEventListener('click', () => {
    fetch('/api/logout', {
      method: 'POST',
      credentials: 'include'
    }).then(res => {
      if (res.ok) {
        logindiv();
        login()
        // connectedUsers.set(res.id, 'offline');
      } else {
        console.log("Logout failed");
      }
    });
  });

}

export function logoutTheUser() {
  let logoutJson = { "ws": ws, "type": "offline" }
  const logouT = JSON.stringify(logoutJson);
  ws.send(logouT)

}