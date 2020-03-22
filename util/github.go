package util

import (
	"context"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v30/github"
	"github.com/pkg/errors"
)

type GitHubClient struct {
	*github.Client
	Token string
}

func NewGitHubClient(appID, installationID int64, privateKey []byte) (*GitHubClient, error) {
	tr, err := ghinstallation.New(http.DefaultTransport, appID, installationID, privateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	token, err := tr.Token(context.TODO())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cli := github.NewClient(&http.Client{Transport: tr})

	return &GitHubClient{
		Client: cli,
		Token:  token,
	}, nil
}
