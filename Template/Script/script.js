function LikePost(postId) {
    fetch(`http://localhost:8666/like/${postId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            likeDislike(data);
        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
}

function DislikePost(postId) {
    fetch(`http://localhost:8666/dislike/${postId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            likeDislike(data);
        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
}

function likeDislike(data) {
    const postsContainer = document.getElementById('post');
    postsContainer.innerHTML = '';
    data.posts.forEach(post => {
        const postDiv = document.createElement('div');
        postDiv.classList.add('post');

        postDiv.innerHTML = `
            <div class="post-header">
                <span class="post-author">${post.author}</span>
                <span class="post-date">${post.created_at}</span>
            </div>
            <h3 class="post-title">${post.title}</h3>
            <p class="post-content">${post.content}</p>
            <div class="category-tags">
                ${post.categories.forEach(cat => `<span class="category-tag">#${cat}</span>`)}
            </div>
            <div class="post-actions">
                <span>Likes: ${post.total_likes}</span>
                <span>Dislikes: ${post.total_dislikes}</span>
                <span>${post.total_comments} Comment(s)</span>
                <button onclick="LikePost(${post.id})">Like</button>
                <button onclick="DislikePost(${post.id})">Dislike</button>
            </div>
        `;

        postsContainer.appendChild(postDiv);
    });
}