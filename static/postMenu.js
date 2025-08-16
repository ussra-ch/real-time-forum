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

export function editpost(postId, postTitle, postContent) {
    const newTitle = prompt("Title:", postTitle);
    const newContent = prompt("content:", postContent);

    if (newTitle === null || newContent === null) return;

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
