import { comment } from "./comment.js";
import { fetchComments } from "./comment.js";
import { deletepost, editpost } from "./postMenu.js";
import { main } from "./main.js";
import { isAuthenticated } from "./login.js";
import { triggerUserLogout } from "./logout.js";

export function categories() {
  let categories = ['All', 'Music', 'Sport', 'Technology', 'Science', 'Culture'];
  const categoDiv = document.getElementById('category');

  categories.forEach(element => {
    const button = document.createElement('button');
    button.className = 'categories';
    button.innerText = element;
    button.addEventListener('click', (e) => {
      isAuthenticated().then(auth => {
        if (!auth) {
          triggerUserLogout()
          main()
        } else {
          e.preventDefault()
          fetch(`/api/fetch_posts`)
            .then(res => res.json())
            .then(posts => {
              const postsContainer = document.getElementById('postsContainer');
              if (!posts) {
                return
              }

              postsContainer.innerHTML = '';
              posts.forEach(post => {
                const topics = post.interest ? post.interest.split(',') : [];
                if (post.interest.split(',').includes(element) || element === 'All') {
                  const postCard = document.createElement('div');
                  postCard.className = 'post-card1';
                  postCard.innerHTML = `
                             <h3>${post.title}</h3>
                             <p>${post.content}</p>
                             <p>Topics: ${topics.join(', ')}</p>
                             <p>Posted by: #${post.nickname || "Unknown"} on ${new Date(post.created_at).toLocaleDateString()}</p>
                              <form class="commentForm">
                                <div class="inputWrapper">
                              <input type="hidden" name="post_id" value="${post.id}">
                               <input type="text" name="content" class="commentInput" placeholder="Write a comment..." required>
                               <button type="submit" class="commentButton"><i class="fa-solid fa-comment"></i></button>
                               <button type="button" class="show"><i class="fa-solid fa-arrow-down"></i></button>
                               </div>
                             </form>
                         `;
                  const menu = document.createElement('div')
                  menu.style.display = 'none'
                  menu.className = 'menu'
                  postCard.prepend(menu)
                  if (post.myId == post.user_id) {
                    const select = document.createElement('button')
                    select.innerHTML = '<i class="fa-solid fa-ellipsis-vertical"></i>'
                    select.className = 'select'
                    postCard.prepend(select)
                    select.addEventListener('click', (e) => {
                      e.preventDefault()
                      menu.style.display = menu.style.display === 'block' ? 'none' : 'block';
                    })
                    const button = document.createElement('button')
                    button.innerHTML = `<i class="fa-solid fa-trash"></i> Delete`
                    menu.prepend(button)
                    button.addEventListener('click', (e) => {
                      e.preventDefault()
                      deletepost(post.id)
                    })
                    const editPost = document.createElement('button')
                    editPost.innerHTML = `<i class="fa-solid fa-file-pen"></i>  Edit `
                    menu.prepend(editPost)
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
                  let offset = 0;
                  const limit = 10;
                  document.querySelector('.show').addEventListener('click', (e) => {
                    e.preventDefault()
                    if (div.style.display === 'none') {
                      div.style.display = 'block'
                      div.innerHTML = ``

                    } else {
                      div.style.display = 'none'
                    }
                    fetchComments(post.id, div, offset, limit);
                  });


                  let lastCall = 0;
                  const delay = 500;

                  div.addEventListener("scroll", () => {
                    if (div.scrollTop + div.clientHeight >= div.scrollHeight) {
                      const now = Date.now();
                      const canCall = now - lastCall >= delay;

                      if (canCall) {
                        lastCall = now;
                        offset += limit;
                        fetchComments(post.id, div, offset, limit);
                      }
                    }
                  });
                  comment(div)
                }
              });

            })
            .catch(err => console.error('Error fetching posts:', err));

        }
      })

    });
    categoDiv.appendChild(button);
  });
}