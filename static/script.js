import { logindiv } from "./login.js";
import { login } from "./login.js";
import { logout } from "./logout.js"
import { Create } from "./post.js"
import { fetchPosts } from "./post.js";
import { catigories } from "./sort.js";
import { comment } from "./comment.js";
import { fetchUser } from "./users.js";
function main() {
  logindiv();
  login()

}
main()
 export function islogin() {

  logout();
  Create();
  fetchPosts();
  catigories();
  comment();
  fetchUser()
}