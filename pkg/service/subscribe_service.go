package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v30/github"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/ryota-sakamoto/github-apps-substrate/pkg/model/commit"
	"github.com/ryota-sakamoto/github-apps-substrate/pkg/model/pull"
	"github.com/ryota-sakamoto/github-apps-substrate/pkg/repository"
)

type repositoryDetail struct {
	owner    string
	name     string
	cloneURL string
}

var (
	appRepositoryURL repositoryDetail = repositoryDetail{
		owner:    "ryota-sakamoto",
		name:     "app-repository",
		cloneURL: "https://github.com/ryota-sakamoto/app-repository.git",
	}
	manifestRepositoryURL repositoryDetail = repositoryDetail{
		owner:    "ryota-sakamoto",
		name:     "manifest-repository",
		cloneURL: "https://github.com/ryota-sakamoto/manifest-repository.git",
	}
)

type SubscribeService interface {
	SubscribePush(event *github.PushEvent) error
}

func NewSubscribeService(repositoryRepository repository.RepositoryRepository) SubscribeService {
	return subscribeService{
		repositoryRepository: repositoryRepository,
	}
}

type subscribeService struct {
	repositoryRepository repository.RepositoryRepository
}

func (s subscribeService) SubscribePush(event *github.PushEvent) error {
	if *event.Deleted {
		return nil
	}

	if event.Repo.Owner.GetName() != appRepositoryURL.owner || event.Repo.GetName() != appRepositoryURL.name {
		log.Printf("push from %s\n", event.Repo.GetURL())
		return nil
	}

	err := s.repositoryRepository.UpdateCommitStatus(context.TODO(), event.Installation.GetID(), commit.UpdateStatus{
		CommitID:    event.HeadCommit.GetID(),
		OwnerName:   appRepositoryURL.owner,
		RepoName:    appRepositoryURL.name,
		Label:       "GitHub Apps",
		Description: "wait",
		Status:      commit.COMMIT_STATUS_PENDING,
	})
	if err != nil {
		return err
	}

	repo, path, err := s.repositoryRepository.CloneRepository(context.TODO(), event.Installation.GetID(), manifestRepositoryURL.cloneURL, event.GetRef())
	if err != nil {
		return err
	}
	defer os.RemoveAll(path)

	current := plumbing.ReferenceName(event.GetRef()).Short()
	log.Println("current:", current)

	if current == "develop" {
		w, _ := repo.Worktree()
		newBranch := fmt.Sprintf("update-develop-%d", time.Now().Unix())
		if err := w.Checkout(&git.CheckoutOptions{
			Create: true,
			Branch: plumbing.NewBranchReferenceName(newBranch),
		}); err != nil {
			return err
		}

		f, err := os.Create(path + "/commit-hash")
		if err != nil {
			return err
		}
		f.Write([]byte(event.HeadCommit.GetID()))
		f.Close()

		w.Add(".")
		w.Commit("update commit hash", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "github-actions[bot]",
				Email: "41898282+github-actions[bot]@users.noreply.github.com",
			},
		})

		if err := repo.Push(&git.PushOptions{}); err != nil {
			log.Printf("%+v\n", err)
			return err
		}

		err = s.repositoryRepository.CreatePullRequest(context.TODO(), event.Installation.GetID(), pull.Request{
			Title: "update commit hash",
			Base:  "develop",
			Head:  newBranch,
			Body:  "update",
		})
		if err != nil {
			log.Printf("%+v\n", err)
			return err
		}
	}

	err = s.repositoryRepository.UpdateCommitStatus(context.TODO(), event.Installation.GetID(), commit.UpdateStatus{
		CommitID:    event.HeadCommit.GetID(),
		OwnerName:   appRepositoryURL.owner,
		RepoName:    appRepositoryURL.name,
		Label:       "GitHub Apps",
		Description: "ok",
		Status:      commit.COMMIT_STATUS_SUCCESS,
	})
	if err != nil {
		return err
	}

	return nil
}
