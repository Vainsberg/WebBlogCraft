package response

type Posts struct {
	Content string
	PostId  string
	Likes   int
}

type PageData struct {
	CurrentPage int
	TotalPages  []int
}

type TemplateData struct {
	Posts      []Posts
	Pagination PageData
}

type PostsRedis struct {
	Content  []string
	PostId   []string
	Likes    []int
	Template TemplateData
}

type PostsIdRedis struct {
	PostId []string
}

type LikeResponse struct {
	NewLikesCount int `json:"newLikesCount"`
}
