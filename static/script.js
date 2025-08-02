import { logindiv } from "./login.js";
import { login } from "./login.js";
import { logout } from "./logout.js"
import { Create } from "./post.js"
import { fetchPosts } from "./post.js";
import { catigories } from "./sort.js";
import { comment } from "./comment.js";
import { initWebSocket } from "./websocket.js";

function main() {
  logindiv();
  login()

}
main()
export function islogin() {
  initWebSocket((msg) => {
    let chatBody = document.getElementById('chat-body');
    let newMsg = document.createElement('div');
    newMsg.innerHTML = `<h3>${msg}</h3>`;
    chatBody.append(newMsg);
  });
  logout();
  Create();
  fetchPosts();
  catigories();
  comment();
}