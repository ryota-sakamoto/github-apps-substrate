package service

import (
	"context"
	"net/http"
	"net/url"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v30/github"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type SubscribeService interface {
	SubscribePush(event *github.PushEvent) error
}

func NewSubscribeService(privateKey string, appID int64) SubscribeService {
	return subscribeService{
		privateKey: privateKey,
		appID:      appID,
	}
}

type subscribeService struct {
	privateKey string
	appID      int64
}

func (s subscribeService) SubscribePush(event *github.PushEvent) error {
	tr, err := ghinstallation.New(http.DefaultTransport, s.appID, event.Installation.GetID(), []byte(s.privateKey))
	if err != nil {
		return errors.WithStack(err)
	}

	token, err := tr.Token(context.TODO())
	if err != nil {
		return errors.WithStack(err)
	}

	u, err := url.Parse(event.Repo.GetCloneURL())
	if err != nil {
		return errors.WithStack(err)
	}

	u.User = url.UserPassword("x-access-token", token)

	fs := memfs.New()
	_, err = git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:           u.String(),
		ReferenceName: plumbing.ReferenceName(event.GetRef()),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
