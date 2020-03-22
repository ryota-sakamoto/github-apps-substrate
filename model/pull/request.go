package pull

type Request struct {
	OwnerName string
	RepoName  string

	Title string
	Head  string
	Base  string
	Body  string
}
