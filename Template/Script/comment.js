document.getElementById('comment-form').addEventListener('submit', function (event) {
    event.preventDefault(); // Prevent form from refreshing the page
    const content = document.getElementById('comment-content').value;

    fetch('/post/{{.Post.ID }}', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({
            'content': content
        })
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json(); // Assuming the server returns the new comment as JSON
        })
        .then(comment => {
            // Append the new comment to the comments section
            const commentsSection = document.getElementById('comments-section');
            const newCommentDiv = document.createElement('div');
            newCommentDiv.className = 'comment';
            newCommentDiv.id = 'comment-' + comment.ID;
            newCommentDiv.innerHTML = `<strong>${comment.Author}</strong> <span>${comment.CreatedAt}</span>
    <p>${comment.Content}</p>
    <span>${comment.TotalLikes} Likes</span> <span>${comment.TotalDislikes} Dislikes</span>`;
            commentsSection.appendChild(newCommentDiv);

            // Clear the textarea
            document.getElementById('comment-content').value = '';
        })
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
        });
});
