package repository

import (
	"context"
	"io/ioutil"
	"net/url"

	"github.com/google/go-github/v30/github"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/ryota-sakamoto/github-apps-substrate/model/commit"
	"github.com/ryota-sakamoto/github-apps-substrate/util"
)

type RepositoryRepository interface {
	CloneRepository(ctx context.Context, installationID int64, repositoryURL, ref string) (string, error)
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

func (r repositoryRepository) UpdateCommitStatus(ctx context.Context, installationID int64, us commit.UpdateStatus) error {
	cli, err := util.NewGitHubClient(r.appID, installationID, r.privateKey)
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

func (r repositoryRepository) CloneRepository(ctx context.Context, installationID int64, repositoryURL, ref string) (string, error) {
	cli, err := util.NewGitHubClient(r.appID, installationID, r.privateKey)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(repositoryURL)
	if err != nil {
		return "", errors.WithStack(err)
	}
	u.User = url.UserPassword("x-access-token", cli.Token)

	path, err := ioutil.TempDir("", "")
	if err != nil {
		return "", errors.WithStack(err)
	}

	_, err = git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
		URL:           u.String(),
		ReferenceName: plumbing.ReferenceName(ref),
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	return path, nil
}
