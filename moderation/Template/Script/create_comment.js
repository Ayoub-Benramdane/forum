function updateCommentCount(increment) {
    const commentButton = document.getElementById("update-number-comment");
    const commentCountSpan = commentButton.querySelector("span");
    const currentCount = parseInt(commentCountSpan.textContent.trim(), 10);
    const newCount = increment ? currentCount + 1 : currentCount - 1;
    commentCountSpan.textContent = ` ${newCount}`;
}

document.getElementById("postCommentBtn").addEventListener("click", async function () {
    const content = document.getElementById("content").value.trim();
    if (content === "") {
        alert("Comment cannot be empty.");
        return;
    }
    const response = await fetch(window.location.pathname, {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `content=${encodeURIComponent(content)}`,
    });

    console.log(response.status);
    console.log(response);
    if (response.ok) {
        const newComment = await response.json();
        const postId = window.location.pathname.split("/").pop();
        appendComment(newComment, postId);
        updateCommentCount(true);
        document.getElementById("content").value = "";
        document.getElementById("noCom").style.display = "none"
    } else {
        alert("Failed to post comment. Please try again.");
    }
});

function appendComment(comment, postId) {
    const commentsSection = document.querySelector(".comments-section");
    const commentDiv = document.createElement("div");
    commentDiv.classList.add("comments-list");
    commentDiv.id = `comment-${comment.id}`;
    const commentHeader = `
        <div class="comment-header">
            <span class="comment-author">${comment.author}</span>
            <span class="comment-date">${comment.created_at}</span>
        </div>`;
    const commentOptions = `
        <div class="comment-options">
            <button id="edit-btn-comment-${postId}-${comment.id}" class="edit-btn"
                onclick="editComment(${postId}, ${comment.id}, 'edit_comment')">
                <i class="fas fa-edit"></i>Edit
            </button>
            <button id="delete-btn-comment-${comment.id}" class="delete-btn"
                onclick="deleteComment(${postId}, ${comment.id}, 'delete_comment')">
                <i class="fas fa-trash"></i>Delete
            </button>
        </div>`
    const commentContent = `<p>${comment.content}</p>`;
    const likeDislike = `
        <div class="like-dislike">
                <button id="like-btn-${postId}${comment.id}" class="action-btn"
                onclick="updateLikeDislikeComment(${postId}, ${comment.id}, 'like_comment')">
                        <i class="fas fa-thumbs-up"></i>
                        <span id="like-count-${postId}${comment.id}"> ${comment.totalLikes || 0}</span>
                </button>
                <button id="dislike-btn-${postId}${comment.id}" class="action-btn"
                onclick="updateLikeDislikeComment(${postId}, ${comment.id}, 'dislike_comment')">
                        <i class="fas fa-thumbs-down"></i>
                        <span id="dislike-count-${postId}${comment.id}"> ${comment.totalLikes || 0}</span>
                </button>
        </div>`;
    commentDiv.innerHTML = commentHeader + commentOptions + commentContent + likeDislike;
    commentsSection.prepend(commentDiv);
}