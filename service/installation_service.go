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
	default:
		return errors.New("invalid action")
	}

	return nil
}
