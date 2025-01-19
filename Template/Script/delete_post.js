async function deletePost(postID, action) {
    const confirmDelete = confirm("Are you sure you want to delete this post?");
    if (!confirmDelete) return;

    try {
        const response = await fetch(`/${action}/${postID}`, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (response.ok) {
            window.location.href = "/";
        } else if (response.status === 400) {
            alert("Invalid post ID.");
        } else if (response.status === 405) {
            alert("Method not allowed.");
        } else if (response.status === 404) {
            alert("post not found. It might have already been deleted.");
        } else {
            alert("You are not authorized to delete this post.");
        }
    } catch (error) {
        alert("An error occurred while trying to delete the post. Please try again.");
        console.error("Delete error:", error);
    }
}