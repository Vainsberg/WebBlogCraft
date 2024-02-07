package response

type PageData struct {
	CurrentPage int
	TotalPages  []int
}

type Comment struct {
	CommentId string
	Comment   string
	UserName  string
	Likes     int
}

type Post struct {
	Content  string
	PostId   string
	UserName string
	Likes    int
	Comments []Comment
}

type ResponseData struct {
	Posts      []Post
	Pagination PageData
}

type LikeResponse struct {
	NewLikesCount int `json:"newLikesCount"`
}

type CommentInput struct {
	Comment string `json:"comment"`
}

type CommentResponse struct {
	CommentID int    `json:"commentId"`
	Comment   string `json:"comment"`
	UserName  string `json:"userName"`
	Likes     int    `json:"likes"`
}

type RabbitMQMessage struct {
	Code  int
	Email string
}
