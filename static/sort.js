import { comment } from "./comment.js";
import { loadComments } from "./comment.js";
export function categories() {
  let categories = ['All', 'Music', 'Sport', 'Tecknology', 'Science', 'Culture'];
  const categoDiv = document.getElementById('catego');

  categories.forEach(element => {
    const boutton = document.createElement('button');
    boutton.className = 'categories';
    boutton.innerText = element;
    boutton.addEventListener('click', (e) => {
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
        ${post.photo ? `<img src="${post.photo}" alt="Post image" style="max-width:100%;">` : ''}
        <p>Posted by: User #${post.nickname || "Unknown"} on ${new Date(post.created_at).toLocaleDateString()}</p>
        <form class="commentForm">
         <input type="hidden" name="post_id" value="${post.id}">
          <input type="text" name="content" class="commentInput" placeholder="Write a comment..." required>
          <button type="submit" class="commentButton">Comment</button>
          <button type="button" class="show">Show Comments</button>
        </form>
      `;

              const div = document.createElement('div');
              div.className = 'comments-container';
              postCard.appendChild(div);
              postsContainer.prepend(postCard);
              document.querySelector('.show').addEventListener('click', (e) => {
                div.style.display = div.style.display === 'none' ? 'block' : 'none';
                loadComments(post.id, div);
              });
              
            }
          });

          comment()
        })
        .catch(err => console.error('Error fetching posts:', err));

    });
    categoDiv.appendChild(boutton);
  });
}