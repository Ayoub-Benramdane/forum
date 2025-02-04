function updateCommentCount(increment) {
    const commentButton = document.getElementById("update-number-comment");
    const commentCountSpan = commentButton.querySelector("span");

    const currentCount = parseInt(commentCountSpan.textContent.trim(), 10);
    const newCount = increment ? currentCount + 1 : currentCount - 1;

    commentCountSpan.textContent = ` ${newCount}`;
}

async function deleteComment(postID, commentID, action) {

    const confirmDelete = confirm("Are you sure you want to delete this comment?");
    if (!confirmDelete) return;

    try {
        const response = await fetch(`/${action}/${postID}/${commentID}`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
        });

        if (response.ok) {
            document.getElementById(`comment-${commentID}`).remove();
            updateCommentCount(false);
        } else if (response.status === 400) {
            alert("Invalid comment ID.");
        } else if (response.status === 405) {
            alert("Method not allowed.");
        } else if (response.status === 404) {
            alert("Comment not found. It might have already been deleted.");
        } else {
            alert("You are not authorized to delete this comment.");
        }
    } catch (error) {
        alert("An error occurred while deleting the comment. Please try again.");
        console.error("Error:", error);
    }
}