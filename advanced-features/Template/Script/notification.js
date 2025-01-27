function markAsRead(notificationId) {
    fetch(`/api/notifications/${notificationId}/mark-read`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => {
            if (response.ok) {
                const notification = document.querySelector(`[data-notification-id="${notificationId}"]`);
                notification.classList.remove('unread');
                notification.querySelector('.mark-read-btn').remove();
                notification.querySelector('.unread-indicator').remove();
            }
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
            } else if (filter === 'unread') {
                notification.style.display = notification.classList.contains('unread') ? 'block' : 'none';
            } else {
                notification.style.display = notification.dataset.type === filter ? 'block' : 'none';
            }
        });
    });
});