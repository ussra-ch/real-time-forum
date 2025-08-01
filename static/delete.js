import { fetchPosts } from "./post.js";

export function deletepost(postId) {
    console.log("Deleting post:", postId);

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
