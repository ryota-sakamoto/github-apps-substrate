package repository

import (
	"context"
	"io/ioutil"
	"net/url"

	"github.com/google/go-github/v30/github"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/ryota-sakamoto/github-apps-substrate/pkg/client"
	"github.com/ryota-sakamoto/github-apps-substrate/pkg/model/commit"
	"github.com/ryota-sakamoto/github-apps-substrate/pkg/model/pull"
)

type RepositoryRepository interface {
	CloneRepository(ctx context.Context, installationID int64, repositoryURL, ref string) (*git.Repository, string, error)
	CreatePullRequest(ctx context.Context, installationID int64, req pull.Request) error
	UpdateCommitStatus(ctx context.Context, installationID int64, us commit.UpdateStatus) error
}

func NewRepositoryRepository(privateKey string, appID int64) RepositoryRepository {
	return repositoryRepository{
		privateKey: []byte(privateKey),
		appID:      appID,
	}
}

type repositoryRepository struct {
	privateKey []byte
	appID      int64
}

func (r repositoryRepository) CloneRepository(ctx context.Context, installationID int64, repositoryURL, ref string) (*git.Repository, string, error) {
	cli, err := client.NewGitHubClient(r.appID, installationID, r.privateKey)
	if err != nil {
		return nil, "", err
	}

	u, err := url.Parse(repositoryURL)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}
	u.User = url.UserPassword("x-access-token", cli.Token)

	path, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	repo, err := git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
		URL:           u.String(),
		ReferenceName: plumbing.ReferenceName(ref),
	})
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	return repo, path, nil
}

func (r repositoryRepository) CreatePullRequest(ctx context.Context, installationID int64, req pull.Request) error {
	cli, err := client.NewGitHubClient(r.appID, installationID, r.privateKey)
	if err != nil {
		return err
	}

	_, _, err = cli.PullRequests.Create(ctx, req.OwnerName, req.RepoName, &github.NewPullRequest{
		Title: &req.Title,
		Base:  &req.Base,
		Head:  &req.Head,
		Body:  &req.Body,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r repositoryRepository) UpdateCommitStatus(ctx context.Context, installationID int64, us commit.UpdateStatus) error {
	cli, err := client.NewGitHubClient(r.appID, installationID, r.privateKey)
	if err != nil {
		return err
	}

	_, _, err = cli.Repositories.CreateStatus(
		ctx,
		us.OwnerName,
		us.RepoName,
		us.CommitID,
		&github.RepoStatus{
			State:       &us.Status,
			Description: &us.Description,
			Context:     &us.Label,
		})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
