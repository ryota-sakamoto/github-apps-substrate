package service

import (
	"context"
	"log"

	"github.com/google/go-github/v30/github"

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

	path, err := s.repositoryRepository.CloneRepository(context.TODO(), event.Installation.GetID(), event.Repo.GetCloneURL(), event.GetRef())
	if err != nil {
		return err
	}

	log.Println(path)

	return nil
}
