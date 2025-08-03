import { fetchUser } from "./users.js";
import { loginDiv, content } from "./var.js";
import { logout } from "./logout.js"
import { Create } from "./post.js"
import { fetchPosts } from "./post.js";
import { catigories } from "./sort.js";
import { comment } from "./comment.js";
import { initWebSocket } from "./websocket.js";


//katjib div dyal login w register (tant que l user ma3andouch session)
export function logindiv() {
    loginDiv.className = 'container';
    loginDiv.id = 'container';
    loginDiv.innerHTML = ` 
  <!-- Sign Up Form -->
  <div class="form-container sign-up-container">
  <form id="reForm" action="/register" method="POST">
  <h1>Create Account</h1>
  <input type="text" placeholder="Nickname" name="Nickname" required />
  <input type="number" placeholder="age" name="Age" min="13" required />
  <select required name="gender">
  <option value="" disabled selected>Gender</option>
  <option value="male">Male</option>
  <option value="female">Female</option>
  </select>
  <input type="text" placeholder="First Name" name="first_name" required />
  <input type="text" placeholder="Last Name" name="last_name" required />
  <input type="email" placeholder="E-mail" name="email" required />
  <input type="password" placeholder="Password"  name="password" required />
  <button id="register">Sign Up</button>
  </form>
  </div>
  <div class="form-container sign-in-container">
  <form id="logForm" action="/login" method="POST">
  <h1>Sign in</h1>
  
  <input type="Nickname" name="Nickname" placeholder="Nickname" />
  <input type="password" name="password" placeholder="Password" />
  <a href="#">Forgot your password?</a>
  <button id="login">Sign In</button>
  </form>
  </div>
  <div class="overlay-container">
  <div class="overlay">
  <div class="overlay-panel overlay-left">
  <h1>Welcome Back!</h1>
  <p>To keep connected with us please login with your personal info</p>
  <button class="ghost" id="signIn">Sign In</button>
  </div>
  <div class="overlay-panel overlay-right">
  <h1>Hello, Friend!</h1>
  <p>Enter your personal details and start journey with us</p>
  <button class="ghost" id="signUp">Sign Up</button>
  </div>
  </div>
  </div>
  
  `;
    document.body.appendChild(loginDiv);
    // Use querySelector on loginDiv to get the buttons
    const signUpButton = loginDiv.querySelector('#signUp');
    const signInButton = loginDiv.querySelector('#signIn');
    const container = loginDiv;
    const form = document.getElementById('logForm')

    document.getElementById('login').addEventListener("click", (e) => {
        e.preventDefault()
        const formData = new FormData(form)
        const data = Object.fromEntries(formData.entries())

        fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        }).then(res => {
            login();
        })
            .catch(err => console.error('Login error:', err));
    })
    const reForm = document.getElementById('reForm')
    document.getElementById('register').addEventListener('click', (e) => {
        e.preventDefault()
        const formData = new FormData(reForm)
        const data = Object.fromEntries(formData.entries())
        fetch('/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        }).then(res => {
            login();
        })
            .catch(err => console.error('Login error:', err));
    })
    signUpButton.addEventListener('click', () => {
        container.classList.add("right-panel-active");
    });

    signInButton.addEventListener('click', () => {
        container.classList.remove("right-panel-active");
    });
}

function islogin() {
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


//check if the user is logged in or not (katchuf session)
export  function login() {
    const body = document.querySelector('body')
    fetch('/api/anthenticated')
        .then(res => {
            if (res.ok) {
                body.innerHTML = `
                <div id="content">
                <header>
                <button id="logout" style="z-index: 10;">log out</button>
                <button id="Create" style="z-index: 10;">+</button>
                </header>
                <div id="catego"></div>
                <div id="all">
             <div id="postsContainer"></div>
            <div id="user">
            <h3>Notifications</h3>
            <div id="not"></div>
            <h3 style="color: rgb(89, 230, 187);"><i class="fa-solid fa-certificate"></i>online</h3>
            <div id="users"></div>
            </div>
            </div>
            </div>
            
            <script type="module" src="static/main.js"></script>
            `
                // fetchUser()
                islogin()
                return true
            } else {
                body.innerHTML = `
    <script type="module" src="static/main.js"></script>
    `
                logindiv()
                return false
            }
        })
}

