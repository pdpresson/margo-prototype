package services

import (
	"context"
	"fmt"
	"gitops_client/appConfig"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-logr/logr"
)

const (
	ORIGIN = "origin"
)

type PullService struct {
	log        *logr.Logger
	appConfig  *appConfig.AppConfig
	appService *AppService
}

func (s PullService) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		sleepDuration := time.Duration(s.appConfig.PollFrequency) * time.Second
		time.Sleep(sleepDuration)

		if s.appConfig.DeviceRepoUrl == "" || s.appConfig.DeviceRepoBranch == "" || s.appConfig.DeviceId == "" {
			s.log.Info("Required configuration not provided. Waiting for configuration...")
			continue
		}

		isNew, err := s.cloneIfNew(ctx)
		if err != nil {
			s.log.Error(err, fmt.Sprint("Error cloning repo", s.appConfig.DeviceRepoUrl))
			continue
		} else if isNew {
			s.appService.ProcessRepo()
			continue
		}

		hasChanges, err := s.checkForChanges(ctx)
		if err != nil {
			s.log.Error(err, fmt.Sprint("Error checking changes for repo", s.appConfig.DeviceRepoUrl))
			continue
		}

		if hasChanges {
			s.log.Info(fmt.Sprint("Changes found in repo ", s.appConfig.DeviceRepoUrl))
			s.appService.ProcessRepo()
		} else {
			s.log.Info(fmt.Sprint("No changes in repo ", s.appConfig.DeviceRepoUrl))
		}

		time.Sleep(time.Duration(s.appConfig.PollFrequency) * time.Second)
	}
}

func (s PullService) cloneIfNew(ctx context.Context) (bool, error) {
	s.log.Info(fmt.Sprintf("Checking if app repo %s is new", s.appConfig.DeviceRepoUrl))

	repoPath := s.appConfig.GetRepoFolder()
	_, err := os.Stat(repoPath)
	if err != nil {
		s.log.Info(fmt.Sprintf("Cloning new repo %s to %s", s.appConfig.DeviceRepoUrl, repoPath))
		repo, err := git.PlainCloneContext(ctx, repoPath, false, &git.CloneOptions{
			URL:               s.appConfig.DeviceRepoUrl,
			Progress:          os.Stdout,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})

		if err != nil {
			return true, err
		}

		err = s.checkoutBranch(repo)
		return true, err
	}

	return false, nil
}

func (s PullService) checkoutBranch(repo *git.Repository) error {
	s.log.Info(fmt.Sprint("Checking out branch ", s.appConfig.DeviceRepoBranch))
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	branchRefName := plumbing.NewBranchReferenceName(s.appConfig.DeviceRepoBranch)
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
		Force:  true,
	}

	err = w.Checkout(&branchCoOpts)
	if err != nil {
		return err
	}

	return nil
}

func (s PullService) checkForChanges(ctx context.Context) (bool, error) {
	s.log.Info("Checking if remote has changes")

	repoPath := s.appConfig.GetRepoFolder()

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return false, err
	}

	hasChanges, err := fetchChanges(repo)
	if err != nil {
		return false, err
	} else if !hasChanges {
		return false, nil
	}

	w, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	s.log.Info(fmt.Sprint("Pulling changes from repo ", repoPath))
	err = w.PullContext(ctx, &git.PullOptions{})
	if err != nil {
		return true, err
	}

	return true, nil
}

func fetchChanges(repo *git.Repository) (bool, error) {
	remote, err := repo.Remote(ORIGIN)
	if err != nil {
		return false, err
	}

	err = remote.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*"},
	})

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func NewPullService(l *logr.Logger, c *appConfig.AppConfig, s *AppService) *PullService {
	return &PullService{
		log:        l,
		appConfig:  c,
		appService: s,
	}
}
