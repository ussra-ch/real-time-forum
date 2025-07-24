const loginDiv = document.createElement('div');
const Logout = document.getElementById('logout');
function logindiv() {
  loginDiv.className = 'container';
  loginDiv.id = 'container';
  loginDiv.innerHTML = ` 
    <!-- Sign Up Form -->
    <div class="form-container sign-up-container">
      <form action="/register" method="POST">
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
        <button>Sign Up</button>
      </form>
    </div>
    <div class="form-container sign-in-container">
      <form action="/login" method="POST">
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
  document.getElementById('login').addEventListener("click", (e) => {
    //e.preventDefault()
  })
  signUpButton.addEventListener('click', () => {
    container.classList.add("right-panel-active");
  });

  signInButton.addEventListener('click', () => {
    container.classList.remove("right-panel-active");
  });
}
logindiv();
function login() {

  fetch('/api/anthenticated')
    .then(res => {
      if (res.ok) {
        loginDiv.style.display = 'none'
        Logout.style.display = 'block'
      } else {
        loginDiv.style.display = 'block'
        Logout.style.display = 'none'
      }
    })
}
login()
function logout() {
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

logout();
