package response

type StoragePostsRedis struct {
	PostsID []string
	Content []string
}

type Posts struct {
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
