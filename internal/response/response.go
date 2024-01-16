package response

type Posts struct {
	Content []string
}
type PostsRedis struct {
	Content []string
}

type PageData struct {
	CurrentPage int
	TotalPages  int
}

type TemplateData struct {
	Posts      []Posts
	Pagination PageData
}
