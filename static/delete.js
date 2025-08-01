function deletepost(postId) {
    fetch('/delete', {
        method: 'POST',
        id: postId,
    }).then(fetchPosts())
}