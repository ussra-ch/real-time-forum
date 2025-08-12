import { comment } from "./comment.js";
import { loadComments } from "./comment.js";
import { deletepost, editpost } from "./postMenu.js";
window.deletepost = deletepost;


export function Create() {
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
  <div class="div-photo">
  <label for="photo">Upload Photo:</label>
  <input type="file" id="photo" name="photo" accept="image/*">
  </div>
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
        CreateCard.style.display = CreateCard.style.display === 'none' ? 'block' : 'none';
    });
    document.getElementById('createPostForm').addEventListener('submit', async function (e) {
        e.preventDefault();
        const selectedTopics = Array.from(document.querySelectorAll('input[name="topic"]:checked')).map(el => el.value);
        const formData = new FormData();
        formData.append('title', this.title.value);
        formData.append('description', this.description.value);
        selectedTopics.forEach(topic => formData.append('topics', topic));
        if (this.photo.files[0]) {
            formData.append('photo', this.photo.files[0]);
        }

        fetch('/api/post', {
            method: 'POST',
            body: formData
        })
            .then(r => { 
                if (!r.ok) {
                return r.json().then(errorData => {
                throw new Error(errorData.Text || `HTTP error! Status: ${r.status}`);
        });
    }
                return r.json();})
            .then(data => {
                console.log("data content is :", data);
                if (data.Type === 'error'){
                        // const existingError = document.querySelector(".errorDiv");
                        // if (existingError) {
                        //     existingError.remove();
                        // }
                        const ErrorDiv = document.createElement('div');
                        ErrorDiv.className = 'error-container';
                        ErrorDiv.innerHTML = `
                                <div class="errorDiv">
                                ${data.Text}
                                </div>`
                        document.querySelector('body').append(ErrorDiv)                        
                    }
                    CreateCard.style.display = 'none';
                    fetchPosts();
                })
            .catch(err => {
                console.error('Error', err.message);
                //createPostCard
                const PostCard = document.getElementById('createPostCard')
                if (PostCard) {
                    PostCard.remove();
                }
                const existingPopup = document.querySelector(".content");
                if (existingPopup) {
                    existingPopup.remove();
                }
                const ErrorDiv = document.createElement('div');
                ErrorDiv.className = 'error-container';
                ErrorDiv.innerHTML = `<div class="content">${err.message}</div>`;
                document.querySelector('body').append(ErrorDiv);
            });
    });


}

export function Notifications(notifs){
    const notifications = document.getElementById('notifications')
    const CreateCard = document.createElement('div')
    CreateCard.innerHTML = `
  <div  id="notifications" > ${notifs}
  </div>
  `;
    content.appendChild(notifications);
    notifications.style.display = 'none';
}


export function fetchPosts() {
    fetch('/api/fetch_posts')
        .then(res => res.json())
        .then(posts => {
            const postsContainer = document.getElementById('postsContainer');
            postsContainer.innerHTML = '';
            if (!posts) {
                return
            }
            posts.forEach(post => {
                const topics = post.interest ? post.interest.split(',') : [];
                const postCard = document.createElement('div');
                postCard.className = 'post-card1';
                postCard.innerHTML = `
                    <h3>${post.title}</h3>
                    <p>${post.content}</p>
                    <p>Topics: ${topics.join(', ')}</p>
                    ${post.photo ? `<img src="${post.photo}" alt="Post image" style="max-width:100%;">` : ''}
                    <p>Posted by: User #${post.nickname || "Unknown"} on ${new Date(post.created_at).toLocaleDateString()}</p>
                     <form class="commentForm">
                     <input type="hidden" name="post_id" value="${post.id}">
                      <input type="text" name="content" class="commentInput" placeholder="Write a comment..." required>
                      <button type="submit" class="commentButton">Comment</button>
                      <button type="button" class="show">Show Comments</button>
                    </form>
                `;
                if (post.myId == post.user_id) {
                    const button = document.createElement('button')
                    button.textContent = 'Delete'
                    postCard.prepend(button)
                    button.addEventListener('click', (e) => {
                        e.preventDefault()
                        deletepost(post.id)
                    })
                    const ed = document.createElement('button')
                    ed.textContent = 'Edit'
                    postCard.prepend(ed)
                    ed.addEventListener('click', (e) => {
                        e.preventDefault()
                        editpost(post.id, post.title, post.content)
                    })
                }



                const div = document.createElement('div');
                div.className = 'comments-container';
                postCard.appendChild(div);
                postsContainer.prepend(postCard);
                div.style.display = 'none'
                document.querySelector('.show').addEventListener('click', (e) => {
                    e.preventDefault()
                    div.style.display = div.style.display === 'none' ? 'block' : 'none';

                    loadComments(post.id, div);
                });
            });
            comment()
        })
        .catch(err => console.error('Error fetching posts:', err));

}