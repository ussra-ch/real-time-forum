import { login } from "./login.js";

export function logout() {
  const Logout = document.getElementById('logout');
  Logout.addEventListener('click', () => {
    fetch('/api/logout', {
      method: 'POST',
      credentials: 'include'
    }).then(res => {
      if (res.ok) {
        login()
      } else {
        console.log("Logout failed");
      }
    });
  });
}