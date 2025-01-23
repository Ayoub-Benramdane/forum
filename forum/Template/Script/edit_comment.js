async function editComment(postID, commentID, action) {

    const confirmDelete = confirm("Are you sure you want to edit this comment?");
    if (!confirmDelete) return;

    try {
        const response = await fetch(`/${action}/${postID}/${commentID}`, {
            method: "GET",
            headers: { "Content-Type": "application/json" },
        });

        if (response.ok) {
            window.location.href = `/${action}/${postID}/${commentID}`;
        } else if (response.status === 400) {
            alert("Invalid comment ID.");
        } else if (response.status === 405) {
            alert("Method not allowed.");
        } else if (response.status === 404) {
            alert("Comment not found. It might have already been edit.");
        } else if (response.status === 401) {
            alert("You are not authorized to edit this comment.");
        } else {
            alert("Oops error in servire.");
        }
    } catch (error) {
        alert("An error occurred while deleting the comment. Please try again.");
        console.error("Error:", error);
    }
}