import { fetchPosts } from "./post.js";

export function deletepost(postId) {

    fetch('/delete', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id: postId }),
    })
        .then(res => {
            if (!res.ok) {
                throw new Error("Failed to delete post");
            }
            return res.text();
        })
        .then(() => {
            fetchPosts();
        })
        .catch(err => {
            console.error("Error deleting post:", err);
        });
}

function showEditModal(postId, postTitle, postContent) {
    let modal = document.getElementById('editPostModal');
    if (!modal) {
        modal = document.createElement('div');
        modal.id = 'editPostModal';
        modal.className = 'edit-modal';
        modal.innerHTML = `
            <div class="edit-modal-content">
                <h2 class="edit-modal-title">Edit Post</h2>
                <input id="editTitle" type="text" placeholder="Title" class="edit-modal-input" />
                <textarea id="editContent" placeholder="Content" class="edit-modal-textarea"></textarea>
                <div class="edit-modal-actions">
                  <button id="saveEditBtn" class="edit-modal-save">Save</button>
                  <button id="cancelEditBtn" class="edit-modal-cancel">Cancel</button>
                </div>
            </div>
        `;
        document.body.appendChild(modal);
    } else {
        modal.style.display = 'flex';
    }

    // Set current values
    modal.querySelector('#editTitle').value = postTitle;
    modal.querySelector('#editContent').value = postContent;

    // Save handler
    modal.querySelector('#saveEditBtn').onclick = () => {
        const newTitle = modal.querySelector('#editTitle').value;
        const newContent = modal.querySelector('#editContent').value;
        modal.style.display = 'none';
        submitEdit(postId, newTitle, newContent);
    };

    // Cancel handler
    modal.querySelector('#cancelEditBtn').onclick = () => {
        modal.style.display = 'none';
    };
}

// Extracted fetch logic
function submitEdit(postId, newTitle, newContent) {
    if (newTitle === "" || newContent === "") return;

    fetch('/edit', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            id: postId,
            title: newTitle,
            content: newContent
        }),
    })
        .then(res => {
            if (!res.ok) {
                throw new Error("error");
            }
            return res.text();
        })
        .then(() => {
            fetchPosts();
        })
        .catch(err => {
            console.error("error", err);
        });
}

export function editpost(postId, postTitle, postContent) {
    showEditModal(postId, postTitle, postContent);
}
