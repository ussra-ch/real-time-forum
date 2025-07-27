
import { comment } from "./comment.js";
import { loadComments } from "./comment.js";
export function catigories() {
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
            console.log(post.interest);

            if (post.interest.split(',').includes(element) || element === 'All') {
              const topics = post.interest ? post.interest.split(',') : [];

              const postCard = document.createElement('div');
              postCard.className = 'post-card1';
              postCard.innerHTML = `
              <h3>${post.title}</h3>
              <p>${post.content}</p>
              <p>Topics: ${topics.join(', ')}</p>
              <p>Posted by: User #${post.Name} on ${new Date(post.created_at).toLocaleDateString()}</p>
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