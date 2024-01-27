package dto

type CommentDto struct {
	CommentId string
	Comment   string
	UserName  string
	Likes     int
}

type PostDto struct {
	Content  string
	PostId   string
	UserName string
	Likes    int
	Comments []CommentDto
}
