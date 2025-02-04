function markAsRead(notificationId) {

    fetch(`/notifications/${notificationId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        
        .catch(error => {
            console.error('Error marking notification as read:', error);
        });
}


function updateFilterCounts() {
    const notifications = document.querySelectorAll('.notification-card');
    const unreadCount = document.querySelectorAll('.notification-card.unread').length;
    const commentCount = document.querySelectorAll('.notification-card[data-type="comment"]').length;
    const reactionCount = document.querySelectorAll('.notification-card[data-type="like"], .notification-card[data-type="dislike"]').length;

    document.querySelectorAll('.filter-btn').forEach(button => {
        const filter = button.dataset.filter;
        let count = 0;

        switch (filter) {
            case 'all':
                count = notifications.length;
                break;
            case 'comment':
                count = commentCount;
                break;
            case 'like':
                count = reactionCount;
                break;
        }

        button.textContent = `${button.textContent.split(' (')[0]} (${count})`;
    });
}

document.querySelectorAll('.filter-btn').forEach(button => {
    button.addEventListener('click', () => {
        document.querySelector('.filter-btn.active').classList.remove('active');
        button.classList.add('active');

        const filter = button.dataset.filter;
        const notifications = document.querySelectorAll('.notification-card');

        notifications.forEach(notification => {
            if (filter === 'all') {
                notification.style.display = 'block';
            } else {
                if (notification.dataset.type == 'comment' && filter == 'comment') {
                    notification.style.display = 'block'
                } else if (notification.dataset.type != 'comment' && filter != 'comment') {
                    notification.style.display = 'block';
                } else {
                    notification.style.display = 'none';
                }
            }
        });
    });
});

document.addEventListener('DOMContentLoaded', () => {
    updateFilterCounts();
});

document.querySelectorAll('.notification-link').forEach(link => {
    link.addEventListener('click', function (event) {
        const notificationCard = this.closest('.notification-card');
        if (notificationCard && notificationCard.classList.contains('unread')) {
            const notificationId = notificationCard.dataset.notificationId;
            markAsRead(notificationId);
        }
    });
});