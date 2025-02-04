async function updateLikeDislike(postID, action) {
    try {
        const response = await fetch(`/${action}/${postID}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error('Failed to update like/dislike');
        }

        const data = await response.json();

        const likeCountElement = document.getElementById(`count-like-${postID}`);
        const dislikeCountElement = document.getElementById(`count-dislike-${postID}`);

        if (data.updatedLikes !== undefined) {
            likeCountElement.innerText = data.updatedLikes;
        }
        if (data.updatedDislikes !== undefined) {
            dislikeCountElement.innerText = data.updatedDislikes;
        }
    } catch (error) {
        console.error('Error updating like/dislike:', error);
    }
}