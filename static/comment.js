export function comment() {
    const forms = document.querySelectorAll('.commentForm');

    forms.forEach((form) => {
        // console.log(1);
        form.addEventListener("submit", (e) => {
            e.preventDefault()

            const commentInput = form.querySelector(".commentInput");
            const comment = commentInput.value;

            const post_id = form.querySelector("[name='post_id']").value;
            if (!comment) return;
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
                    loadComments(post_id,forms)
                })
                .catch(err => {
                    console.error("Error:", err);
                });
        });
    });
}

//fetch comments (l kola post)
export function loadComments(postId, container) {
   container.innerHTML = ``
    fetch('/api/fetch_comments')
        .then(res => res.json())
        .then(comments => {
            if (!comments) {
                return
            }
            comments.forEach(comment => {
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