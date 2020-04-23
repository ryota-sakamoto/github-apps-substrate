package client

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

func NewGitHubClient(baseURL, uploadURL string, appID, installationID int64, privateKey []byte) (*GitHubClient, error) {
	tr, err := ghinstallation.New(http.DefaultTransport, appID, installationID, privateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if baseURL != "" && uploadURL != "" {
		tr.BaseURL = baseURL
	}

	token, err := tr.Token(context.TODO())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var cli *github.Client
	if baseURL != "" && uploadURL != "" {
		cli, err = github.NewEnterpriseClient(baseURL, uploadURL, &http.Client{Transport: tr})
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		cli = github.NewClient(&http.Client{Transport: tr})
	}

	return &GitHubClient{
		Client: cli,
		Token:  token,
	}, nil
}
