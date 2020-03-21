package service

import (
	"context"
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v30/github"
	"github.com/pkg/errors"
)

type InstallationService interface {
	Action(event *github.InstallationEvent) error
}

func NewInstallationService(privateKey string) InstallationService {
	return installationService{
		privateKey: privateKey,
	}
}

type installationService struct {
	privateKey string
}

func (s installationService) Action(event *github.InstallationEvent) error {
	switch event.GetAction() {
	case "created", "added", "removed":
		tr, err := ghinstallation.New(http.DefaultTransport, event.Installation.GetAppID(), event.Installation.GetID(), []byte(s.privateKey))
		if err != nil {
			return errors.WithStack(err)
		}

		cli := github.NewClient(&http.Client{Transport: tr})
		list, _, _ := cli.Apps.ListRepos(context.TODO(), nil)
		for _, v := range list {
			log.Println(v.GetFullName())
		}
	default:
		return errors.New("invalid action")
	}

	return nil
}
