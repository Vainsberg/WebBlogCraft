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

type RandomPostId struct {
	RandPostId int
}

type PostsRedis struct {
	Content  []string
	Template TemplateData
	Random   RandomPostId
}

type PostsIdRedis struct {
	PostId []string
}
