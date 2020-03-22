package commit

type UpdateStatus struct {
	OwnerName string
	RepoName  string
	CommitID  string

	Label       string
	Description string
	Status      string
}
