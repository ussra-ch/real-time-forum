const loginDiv = document.createElement('div');
const content = document.getElementById('content')
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
        content.style.display = 'block'
      } else {
        loginDiv.style.display = 'block'
        content.style.display = 'none'
      }
    })
}
login()
function logout() {
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
logout();
function Create() {
  const Create = document.getElementById('Create')
  const CreateCard = document.createElement('div')
  CreateCard.innerHTML = `
  <div class="post-card" id="createPostCard" >
  <form id="createPostForm" method="get">
  <div class="container">
  <h3>Create Post</h3>
  <div class="div-title">
  <label for="title">Title :</label>
  <input type="text" name="title" id="title" required>
  </div>
  <div class="div-description">
  <label for="description">description :</label>
  <textarea name="description" id="description" rows="4" required></textarea>
  </div>
  <div class="topic-options">
  <label><input type="checkbox" id="music" name="topic" value="Music"> Music</label>
  <label><input type="checkbox" id="sport" name="topic" value="Sport"> Sport</label>
  <label><input type="checkbox" id="gaming" name="topic" value="Gaming"> Gaming</label>
  <label><input type="checkbox" id="health" name="topic" value="Health"> Health</label>
  <label><input type="checkbox" id="general" name="topic" value="General"> General</label>
  </div>
  <div id="errorMsg" style="display:none; color:red; margin: 10px 10px;"></div>
  <button type="submit">Post</button>
  </div>
  </form>
  </div>
  `;
  content.appendChild(CreateCard);
  CreateCard.style.display = 'none';

  Create.addEventListener('click', (e) => {

    CreateCard.style.display = 'block';
  });
  document.getElementById('createPostForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const selectedTopics = Array.from(document.querySelectorAll('input[name="topic"]:checked')).map(el => el.value);

    const data = {
      title: this.title.value,
      description: this.description.value,
      topics: selectedTopics,
    };

    fetch('/api/post', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    })
      .then(r => r.json())
      .then(data => {
        CreateCard.style.display = 'none';
        fetchPosts();
      })
      .catch(err => {
        console.error('Error:', err.message);
      });
  });


}
Create();
function fetchPosts() {
  fetch('/api/fetch_posts')
    .then(res => res.json())
    .then(posts => {
      const postsContainer = document.getElementById('postsContainer');
      postsContainer.innerHTML = '';
      posts.forEach(post => {
        const topics = post.interest ? post.interest.split(',') : [];

        const postCard = document.createElement('div');
        postCard.className = 'post-card1';

        postCard.innerHTML = `
          <h3>${post.title}</h3>
          <p>${post.content}</p>
          <p>Topics: ${topics.join(', ')}</p>
          <p>Posted by: User #${post.user_id} on ${new Date(post.created_at).toLocaleDateString()}</p>
           <form class="commentForm">
            <input type="text" name="post_id" value="${post.id}" hidden>
              <input type="text" name="content" class="commentInput" placeholder="Write a comment..." required>
              <button type="submit" class="commentButton">Comment</button>
              <button  class="show">show Comment</button>
              </form>
          `;
        const div = document.createElement('div');
        div.className = 'comments-container';
        postCard.appendChild(div);
        postsContainer.prepend(postCard);
        document.querySelector('.show').addEventListener('click', (e) => {
          div.style.display = div.style.display === 'none' ? 'block' : 'none';

        });
        loadComments(post.id, div);
        comment();
      });
    })
    .catch(err => console.error('Error fetching posts:', err));
}
fetchPosts();
function catigories() {
  let categories = ['All', 'Music', 'Sport', 'Gaming', 'Health', 'General'];
  const categoDiv = document.getElementById('catego');
  categories.forEach(element => {

    const boutton = document.createElement('button');
    boutton.className = 'catigories';
    boutton.innerText = element;
    boutton.addEventListener('click', () => {
      fetch(`/api/fetch_posts`)
        .then(res => res.json())
        .then(posts => {
          const postsContainer = document.getElementById('postsContainer');
          postsContainer.innerHTML = '';
          posts.forEach(post => {
            if (post.interest == element || element === 'All') {
              const topics = post.interest ? post.interest.split(',') : [];

              const postCard = document.createElement('div');
              postCard.className = 'post-card1';
              postCard.innerHTML = `
              <h3>${post.title}</h3>
              <p>${post.content}</p>
              <p>Topics: ${topics.join(', ')}</p>
              <p>Posted by: User #${post.user_id} on ${new Date(post.created_at).toLocaleDateString()}</p>
              <form class="commentForm">
                  <input type="text" name="post_id" value="${post.id}" hidden>
                  <input type="text" name"content" class="commentInput" placeholder="Write a comment..." required>
                  <button type="submit" class="commentButton">Comment</button>
                  <button  class="show">show Comment</button>
              </form>
              
            `;
              const div = document.createElement('div');
              div.className = 'comments-container';
              postCard.appendChild(div);
              postsContainer.prepend(postCard);
              document.querySelector('.show').addEventListener('click', (e) => {
                div.style.display = div.style.display === 'none' ? 'block' : 'none';

              });
              loadComments(post.id, div);
            }

          });
        })
        .catch(err => console.error('Error fetching posts:', err));
      comment();

    });
    categoDiv.appendChild(boutton);
  });
}
catigories();
function comment() {

  const forms = document.querySelectorAll('.commentForm');

  forms.forEach((form) => {
    form.addEventListener("submit", function (e) {
      e.preventDefault();

      const commentInput = form.querySelector(".commentInput");
      const comment = commentInput.value;

      const post_id = form.querySelector("[name='post_id']").value;

      fetch("/comment", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ comment, post_id }),
      })
        .then(res => res.json())
        .then(data => {
          commentInput.value = "";
        })
        .catch(err => {
          console.error("Error:", err);
        });

    });
  });
}

comment();
function loadComments(postId, container) {

  fetch('/api/fetch_comments')
    .then(res => res.json())
    .then(comments => {
      comments.forEach(comment => {
        console.log(comment.PostID, postId);
        if (comment.PostID != postId) return;

        const p = document.createElement("div");
        p.innerHTML = `
        <p><strong>${comment.Name}:</strong> ${comment.Content}</p>
        <p class="comment-date">${new Date(comment.CreatedAt).toLocaleDateString()}</p>
      `;
        container.appendChild(p);
      });
    });
}
function fetchUser() {
  const usern = document.getElementById('user')
  fetch('/user').then(r => r.json()).then(users => {
    users.forEach(user => {
      const div = document.createElement('div');
      div.innerHTML=`
      <i class="fa-solid fa-user"></i> ${user}
      `
      div.style.border = '1px solid #ccc';
      div.style.padding = '8px';
      div.style.borderRadius = '70px';
      div.style.width='50%'
      div.style.background = ' rgba(26, 35, 50, 0.95)';
      usern.appendChild(div);
    })
  })
}
fetchUser()