package response

type PageData struct {
	CurrentPage int
	TotalPages  []int
}

type Comments struct {
	Comment  string
	UserName string
	Like     string
}

type Posts struct {
	Content  string
	PostId   string
	UserName string
	Likes    int
	Comment  []Comments
}

type TemplateData struct {
	Posts      []Posts
	Pagination PageData
}

type LikeResponse struct {
	NewLikesCount int `json:"newLikesCount"`
}
