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
  <label><input type="checkbox" id="technology" name="topic" value="Technology"> Technology</label>
  <label><input type="checkbox" id="culture" name="topic" value="Culture"> Culture</label>
  <label><input type="checkbox" id="gcience" name="topic" value="Science"> Science</label>
  </div>
  <div id="errorMsg" style="display:none; color:red; margin: 10px 10px;"></div>
  <button type="submit">Post</button>
  </div>
  </form>
  </div>
  `;
    content.appendChild(CreateCard);
    CreateCard.style.display = 'none';

    const deleteButton = document.createElement('button')
    const icon = document.createElement('i');
    icon.className = 'fas fa-trash';
    deleteButton.id = 'closeConversation'
    deleteButton.appendChild(icon);
    deleteButton.appendChild(document.createTextNode('Delete'));
    let createPostDiv = document.getElementById('createPostCard')
    createPostDiv.append(deleteButton)
    // body.append(conversation)

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
                return r.json();
            })
            .then(data => {
                console.log("data content is :", data);
                if (data.Type === 'error') {
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
                // const PostCard = document.getElementById('createPostCard')
                // if (PostCard) {
                //     PostCard.style.display = 'none';
                // }
                const existingPopup = document.querySelector(".content");
                if (existingPopup) {
                    existingPopup.remove();
                }
                const ErrorDiv = document.createElement('div');
                ErrorDiv.className = 'error-container';
                ErrorDiv.innerHTML = `<div class="content">${err.message}</div>`;
                document.querySelector('body').append(ErrorDiv);
                setTimeout(() => {
                    ErrorDiv.remove()
                }, 1000)
            });
    });


    deleteButton.addEventListener('click', () => {
        createPostDiv.remove()
    })



}

export function Notifications(notifs) {
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
                       <div class="inputWrapper">
                     <input type="hidden" name="post_id" value="${post.id}">
                      <input type="text" name="content" class="commentInput" placeholder="Write a comment..." required>
                      <button type="submit" class="commentButton"><i class="fa-solid fa-comment"></i></button>
                      <button type="button" class="show">Show Comments</button>
                      </div>
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
                    const editPost = document.createElement('button')
                    editPost.textContent = 'Edit'
                    postCard.prepend(editPost)
                    editPost.addEventListener('click', (e) => {
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
                comment(div)
            });
        })
        .catch(err => console.error('Error fetching posts:', err));

}