document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.like-button-comment').forEach(button => {
        button.addEventListener('click', function(event) {
            event.preventDefault();
            const commentId = this.getAttribute('comment-id');
            const likesCountElement = document.getElementById('likes-count-comment-' + commentId);
            console.log(commentId);
            fetch(`/posts/${commentId}/comment/like`, {
                method: 'POST',
            })
            .then(response => {
              console.log('Ответ сервера:', response);
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