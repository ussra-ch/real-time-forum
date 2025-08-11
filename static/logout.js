import { login } from "./login.js";


//katloggouti (katfetshi data mn backend bach katms7 session)
export function logout() {
  const Logout = document.getElementById('logout');

  Logout.addEventListener('click', () => {
    fetch('/api/logout', {
      method: 'POST',
      credentials: 'include'
    }).then(res => {
      if (res.ok) {
        login()
        // connectedUsers.set(res.id, 'offline');
      } else {
        console.log("Logout failed");
      }
    });
  });

}