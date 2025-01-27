document.querySelectorAll('.filter-btn').forEach(button => {
    button.addEventListener('click', () => {
        document.querySelector('.filter-btn.active').classList.remove('active');
        button.classList.add('active');

        const filter = button.dataset.filter;
        const activities = document.querySelectorAll('.activity-card');

        activities.forEach(activity => {
            if (filter === 'all') {
                activity.style.display = 'flex';
            } else {
                activity.style.display =
                    (filter === 'posts' && activity.dataset.type === 'post') || (filter === 'comments' && activity.dataset.type === 'comment') ||  (filter === 'reactions' && activity.dataset.type === 'reactions') ? 'flex' : 'none';
            }
        });
    });
});