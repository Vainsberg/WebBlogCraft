document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.like-button').forEach(button => {
        button.addEventListener('click', function(event) {
            event.preventDefault();
            const postId = this.getAttribute('post-id');
            const likesCountElement = document.getElementById('likes-count-' + postId);
            console.log(postId);
            fetch(`/posts/${postId}/like`, {
                method: 'POST',
            })
            .then(response => {
                if (response.ok) {
                    return response.json(); 
                }
                throw new Error('Network response was not ok.');
            })
            .then(data => {
                if (!data.isAuthenticated) {
                    alert('Вы не авторизованы! Пожалуйста, войдите в систему.');
                } 
                if (likesCountElement) {
                    likesCountElement.textContent = data.newLikesCount;
                }
            })
            .catch(error => {
                console.error('Ошибка при добавлении лайка:', error);
            });
        });
    });
});