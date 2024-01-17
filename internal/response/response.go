package response

type Posts struct {
	Content []string
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
	Template TemplateData
}

type PostsIdRedis struct {
	PostId []string
}
