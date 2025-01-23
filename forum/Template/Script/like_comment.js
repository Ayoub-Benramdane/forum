async function updateLikeDislikeComment(PostID, CommentID, action) {
    try {
        const response = await fetch(`/${action}/${PostID}/${CommentID}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error('Failed to update like/dislike');
        }

        const data = await response.json();

        const likeCountElement = document.getElementById(`like-count-${PostID}${CommentID}`);
        const dislikeCountElement = document.getElementById(`dislike-count-${PostID}${CommentID}`);

        if (data.updatedLikes !== undefined) {
            likeCountElement.innerText = data.updatedLikes;
        }
        if (data.updatedDislikes !== undefined) {
            dislikeCountElement.innerText = data.updatedDislikes;
        }

        const likeButton = document.querySelector(`#like-btn-${PostID}${CommentID}`);
        const dislikeButton = document.querySelector(`#dislike-btn-${PostID}${CommentID}`);

        if (data.isLiked) {
            likeButton.classList.add('active');
        } else {
            likeButton.classList.remove('active');
        }

        if (data.isDisliked) {
            dislikeButton.classList.add('active');
        } else {
            dislikeButton.classList.remove('active');
        }
    } catch (error) {
        console.error('Error updating like/dislike:', error);
    }
}