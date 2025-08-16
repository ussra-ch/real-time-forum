import { logindiv } from "./login.js";
import { login } from "./login.js";
import { errorPage } from "./errorPage.js";


export function main() {
  if (document.getElementById('style')) {
    document.getElementById('style').remove()
  }
  const currentUrl = window.location.href;
  const urlArr = currentUrl.split('/')
  console.log(urlArr);
  
  if (urlArr[urlArr.length - 1] != "" || urlArr.length != 4) {
    errorPage('404')
    return
  }

  logindiv();
  login()
}


main()
