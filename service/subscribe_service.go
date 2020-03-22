package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v30/github"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
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

	fs := memfs.New()
	_, err = git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: fmt.Sprintf("https://x-access-token:%s@github.com/owner/repo.git", token),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	file, _ := fs.Open("README.md")
	b, _ := ioutil.ReadAll(file)
	log.Println(string(b))

	return nil
}
