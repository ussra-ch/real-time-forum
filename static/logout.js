import { login, logindiv } from "./login.js";
import { fetchUser } from "./users.js";
import { ws } from "./websocket.js";


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
      } else {
        console.log("Logout failed");
      }
    });
  });

}

export function triggerUserLogout() {
  // let isConversationOpen = {
  //   senderId: userId,
  //   receiverId: 0,
  //   type: "CloseConversation"
  // }

  // ws.send(JSON.stringify(isConversationOpen));
  let logoutJson = { "ws": ws, "type": "offline" }
  const logouT = JSON.stringify(logoutJson);
  ws.send(logouT)
}