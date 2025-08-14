import { logindiv } from "./login.js";
import { login } from "./login.js";
import { Errorr } from "./errorPage.js";


export function main() {
  const currentUrl = window.location.href;
  const urlArr = currentUrl.split('/')

  if (urlArr[urlArr.length - 1] != "" || urlArr.length != 4) {
    Errorr('404')
   
    
    return
  }


  logindiv();
  login()
}


main()
