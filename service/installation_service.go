package service

import (
	"log"

	"github.com/google/go-github/v30/github"
	"github.com/pkg/errors"
	"github.com/ryota-sakamoto/github-apps-substrate/util"
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
		log.Println(util.GetToken(int(*event.Installation.AppID), s.privateKey))
	default:
		return errors.New("invalid action")
	}

	return nil
}
