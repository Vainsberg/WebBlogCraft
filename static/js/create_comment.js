document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.create-comment').forEach(button => {
      button.addEventListener('click', function(event) {
          event.preventDefault();
          const postId = this.getAttribute('post-id');
          const commentInput = document.getElementById('comment-' + postId);
          const commentText = commentInput.value;
          const commentsContainer = document.getElementById('comments-container-' + postId);
  
          fetch(`/posts/${postId}/comment`, {
              method: 'POST',
              headers: {'Content-Type': 'application/json'},
              body: JSON.stringify({ comment: commentText })
          })
          .then(response => response.json())
          .then(data => {
            if (!data.isAuthenticated) {
                alert('Вы не авторизованы! Пожалуйста, войдите в систему.');
            }
              commentInput.value = ''; 
              const newCommentDiv = document.createElement("div");
              newCommentDiv.innerHTML = `
                  <strong>${data.userName}</strong>: ${data.commentId}
                  <button class="like-button-comment" comment-id="${data.commentId}">Лайк</button>
                  <span class="likes-count-comment" id="likes-count-comment-${data.commentId}">${data.likes || 0}</span>
              `;
              commentsContainer.appendChild(newCommentDiv);
  
              newCommentDiv.querySelector('.like-button-comment').addEventListener('click', function(event) {
                  event.preventDefault();
                  const commentId = this.getAttribute('comment-id');
                  const likesCountElement = document.getElementById('likes-count-comment-' + commentId);
  
                  fetch(`/posts/${commentId}/comment//like`, {
                      method: 'GET',
                  })
                  .then(response => response.json())
                  .then(data => {
                    if (!data.isAuthenticated) {
                        alert('Вы не авторизованы! Пожалуйста, войдите в систему.');
                    }
                      if (likesCountElement) {
                          likesCountElement.textContent = data.newLikesCount;
                      }
                  })
                  .catch(error => console.error('Ошибка при добавлении лайка:', error));
              });
          })
          .catch(error => console.error('Ошибка при добавлении комментария:', error));
      });
    });
  });