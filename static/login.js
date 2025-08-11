import { fetchUser } from "./users.js";
import { loginDiv, content, notifications } from "./var.js";
import { profile } from "./profile.js"
import { formatDate } from "./message.js"
import { logout } from "./logout.js"
import { Create } from "./post.js"
import { fetchPosts } from "./post.js";
import { catigories } from "./sort.js";
import { comment } from "./comment.js";
import { initWebSocket } from "./websocket.js";

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
        // console.log(msg);
        let chatBody = document.getElementById('chat-body');
        if (!chatBody || msg == "") {
            return
        }
        let newMsg = document.createElement('div');
        newMsg.innerHTML = `<h3>${msg}</h3>
                            <h7>${formatDate(Date.now())}</h7>`;
        newMsg.className = 'messageSent'
        chatBody.append(newMsg);
        chatBody.scrollTop = chatBody.scrollHeight;
    });
    logout();
    Create();
    fetchPosts();
    catigories();
    comment();
}

export function login() {
    // console.log('dkhal l log function');
    const body = document.querySelector('body')
    fetch('/api/anthenticated')
        .then(r => r.json())
        .then(res => {
            let profil = `<i class="fa-solid fa-user"></i>`
            if (res.photo) {
                profil = `<img src="${res.photo}" alt="Profile Picture">`
            }
            if (res.ok) {
                body.innerHTML = `
                <div id="content">
                <header>
                <button id="profile" style="z-index: 10;">${profil} </button>
                <button id="Create" style="z-index: 10;">+</button>
                <div class="notification-circle">
                    🔔
                    <div class="notification-badge" , id ="notification-circle">${notifications}</div>
                </div>
                </header>
                <div class="sidebar-left">
                <div class="sidebar-label categories-label">Categories</div>
                <div class="sidebar-label posts-label">Posts</div>
                </div>
                <div id="catego"></div>
                <div id="all">
                <div id="postsContainer"></div>
                <div id="user">
                <h3 style="color: rgb(89, 230, 187);"><i class="fa-solid fa-certificate"></i>online</h3>
                <div id="users"></div>
                </div>
                </div>
                </div>
            
            <script type="module" src="static/main.js"></script>
            `
                const div = document.createElement('div');
                div.innerHTML = `
                    <button id="logout">Logout</button>
                    <button id="editProfile">Edit Profile</button>
                `;
                body.append(div);

                div.style.position = 'absolute';
                div.style.top = '8vh';
                div.style.height = '20vh'
                div.style.right = '0';
                div.style.background = 'rgba(26, 35, 50, 0.8)';
                div.style.padding = '10px';
                div.style.boxShadow = '0 2px 8px rgba(0,0,0,0.2)';
                div.style.zIndex = '1000';
                div.style.display = 'none';
                const logoutBtn = document.getElementById('logout');
                const editBtn = document.getElementById('editProfile');

                logoutBtn.style.margin = '5px';
                editBtn.style.position = 'relative';
                editBtn.style.top = '8vh';
                editBtn.style.height = '5vh';

                document.getElementById('profile').addEventListener('click', () => {
                    if (div.style.display === 'none') {
                        div.style.display = 'flex';
                    } else {
                        div.style.display = 'none';
                    }
                });
                document.getElementById('editProfile').addEventListener('click', () => {
                    profile(res.age, res.email, res.nickname, res.photo)
                })
                // console.log(ws.readyState);
                islogin();
                fetchUser(res.status);
                return true
            } else {
                body.innerHTML = `
                    <script type="module" src="static/main.js"></script>
                    `
                logindiv()
                return false
            }
        }).catch(err => console.error('Error:', err));
}

