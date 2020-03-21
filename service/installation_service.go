package service

import (
	"errors"
	"log"

	"github.com/google/go-github/v30/github"
)

type InstallationService interface {
	Action(event *github.InstallationEvent) error
}

func NewInstallationService() InstallationService {
	return installationService{}
}

type installationService struct {
}

func (s installationService) Action(event *github.InstallationEvent) error {
	switch event.GetAction() {
	case "created":
	case "added":
	case "removed":
		log.Println(event)
	default:
		return errors.New("invalid action")
	}

	return nil
}
