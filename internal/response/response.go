package response

type PageData struct {
	CurrentPage int
	TotalPages  []int
}

type Posts struct {
	Content string
	PostId  string
	Likes   int
}

type TemplateData struct {
	Posts      []Posts
	Pagination PageData
}

type LikeResponse struct {
	NewLikesCount int `json:"newLikesCount"`
}
