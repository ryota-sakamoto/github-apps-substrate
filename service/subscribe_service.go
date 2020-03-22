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

	"github.com/ryota-sakamoto/github-apps-substrate/model/commit"
	"github.com/ryota-sakamoto/github-apps-substrate/repository"
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

	err := s.repositoryRepository.UpdateCommitStatus(context.TODO(), event.Installation.GetID(), commit.UpdateStatus{
		CommitID:    event.HeadCommit.GetID(),
		OwnerName:   event.Repo.Owner.GetName(),
		RepoName:    event.Repo.GetName(),
		Label:       "GitHub Apps",
		Description: "wait",
		Status:      commit.COMMIT_STATUS_PENDING,
	})
	if err != nil {
		return err
	}

	go func() {
		repo, path, err := s.repositoryRepository.CloneRepository(context.TODO(), event.Installation.GetID(), event.Repo.GetCloneURL(), event.GetRef())
		if err != nil {
			log.Printf("%+v\n", err)
			return
		}
		defer os.RemoveAll(path)

		current := plumbing.ReferenceName(event.GetRef()).Short()
		log.Println("current:", current)

		if current != "develop" {
			err := s.repositoryRepository.UpdateCommitStatus(context.TODO(), event.Installation.GetID(), commit.UpdateStatus{
				CommitID:    event.HeadCommit.GetID(),
				OwnerName:   event.Repo.Owner.GetName(),
				RepoName:    event.Repo.GetName(),
				Label:       "GitHub Apps",
				Description: "ok",
				Status:      commit.COMMIT_STATUS_SUCCESS,
			})
			if err != nil {
				log.Printf("%+v\n", err)
				return
			}
			return
		}

		log.Println(path)
		w, _ := repo.Worktree()
		if err := w.Checkout(&git.CheckoutOptions{
			Create: true,
			Branch: plumbing.NewBranchReferenceName(fmt.Sprintf("feature/%d", time.Now().UnixNano())),
		}); err != nil {
			log.Printf("%+v\n", err)
			return
		}

		if err := repo.Push(&git.PushOptions{}); err != nil {
			log.Printf("%+v\n", err)
			return
		}
	}()

	return nil
}
