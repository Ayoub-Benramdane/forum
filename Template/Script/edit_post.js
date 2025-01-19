async function editPost(postID, action) {
    const confirmDelete = confirm("Are you sure you want to edit this post?");
    if (!confirmDelete) return;

    try {
        const response = await fetch(`/${action}/${postID}`, {
            method: "GET",
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (response.ok) {
            window.location.href = `/${action}/${postID}`;
        } else if (response.status === 400) {
            alert("Invalid post ID.");
        } else if (response.status === 405) {
            alert("Method not allowed.");
        } else if (response.status === 404) {
            alert("post not found. It might have already been edit.");
        } else if (response.status === 401) {
            alert("You are not authorized to edit this post..");
        } else {
            alert("Oops error in servire.");
        }
    } catch (error) {
        alert("An error occurred while trying to delete the post. Please try again.");
        console.error("Delete error:", error);
    }
}