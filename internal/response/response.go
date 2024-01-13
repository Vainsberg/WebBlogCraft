package response

type Page struct {
	ID    string
	Posts string
}

type StoragePosts struct {
	PostsID []string
	Posts   []string
}

type Posts struct {
	Posts []string
}
