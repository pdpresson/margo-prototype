package services

import (
	"context"
	"fmt"
	"gitops_pullservice/appConfig"
	"gitops_pullservice/clients"
	"os"
	"path/filepath"
	"strings"
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
	log          *logr.Logger
	configClient *clients.PullConfigClient
	appConfig    *appConfig.AppConfig
	appService   *AppService
}

func (s PullService) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		s.log.Info("Updating app repositories...")
		for _, v := range s.configClient.ReposToMonitor() {
			isNew, err := cloneIfNew(ctx, s.log, s.appConfig, v.Url, v.Branch)
			if err != nil {
				s.log.Error(err, fmt.Sprint("Error cloning repo", v.Url))
				continue
			} else if isNew {
				f := getFolderPath(s.appConfig, v.Url)
				s.appService.ProcessRepo(v.ID, f)
				continue
			}

			hasChanges, err := checkForChanges(ctx, s.log, s.appConfig, v.Url)
			if err != nil {
				s.log.Error(err, fmt.Sprint("Error checking changes for repo", v.Url))
				continue
			}

			if hasChanges {
				s.log.Info(fmt.Sprint("Changes found in repo ", v.Url))
				f := getFolderPath(s.appConfig, v.Url)
				s.appService.ProcessRepo(v.ID, f)
			} else {
				s.log.Info(fmt.Sprint("No changes in repo ", v.Url))
			}
		}

		time.Sleep(s.configClient.PollFrequency())
	}
}

func cloneIfNew(ctx context.Context, l *logr.Logger, c *appConfig.AppConfig, repoUrl, branch string) (bool, error) {
	l.Info(fmt.Sprintf("Checking if app repo %s is new", repoUrl))

	repoPath := getFolderPath(c, repoUrl)
	_, err := os.Stat(repoPath)
	if err != nil {
		l.Info(fmt.Sprintf("Cloning new repo %s to %s", repoUrl, repoUrl))
		repo, err := git.PlainCloneContext(ctx, repoPath, false, &git.CloneOptions{
			URL:               repoUrl,
			Progress:          os.Stdout,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})

		if err != nil {
			return true, err
		}

		err = checkoutBranch(l, repo, branch)
		return true, err
	}

	return false, nil
}

func checkoutBranch(l *logr.Logger, repo *git.Repository, branch string) error {
	l.Info(fmt.Sprint("Checking out branch ", branch))
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	branchRefName := plumbing.NewBranchReferenceName(branch)
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

func checkForChanges(ctx context.Context, l *logr.Logger, c *appConfig.AppConfig, repoUrl string) (bool, error) {
	l.Info(fmt.Sprintf("Checking if remote %s has changes", repoUrl))

	repoPath := getFolderPath(c, repoUrl)

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

	l.Info(fmt.Sprint("Pulling changes from repo ", repoPath))
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

func getFolderPath(c *appConfig.AppConfig, repoUrl string) string {
	p := strings.SplitAfter(repoUrl, "/")
	n := p[len(p)-1]
	n = strings.TrimSuffix(n, ".git")

	repoPath := filepath.Join(c.RepoRootPath, n)
	return repoPath
}

/*
func identifyChangedFiles(repo *git.Repository, branch string, lastCheck *time.Time) (map[string]bool, error) {
	changedFiles := make(map[string]bool)

	localRef, err := repo.Reference(plumbing.NewBranchReferenceName(branch), true)
	if err != nil {
		return changedFiles, err
	}

	remoteRef, err := repo.Reference(plumbing.NewRemoteReferenceName(ORIGIN, branch), true)
	if err != nil {
		return changedFiles, err
	}

	if localRef.Hash() == remoteRef.Hash() {
		return changedFiles, err
	}

	fmt.Println("Examining changes since", lastCheck)
	commits, err := repo.Log(&git.LogOptions{
		From:  remoteRef.Hash(),
		Since: lastCheck,
		Order: git.LogOrderCommitterTime,
	})

	if err != nil {
		return changedFiles, err
	}

	err = commits.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		fe, err := c.Files()
		if err != nil {
			return nil
		}
		fe.ForEach(func(f *object.File) error {
			fmt.Printf("File %s changed\n", f.Name)
			changedFiles[f.Name] = true
			return nil
		})
		return nil
	})

	return changedFiles, err
}


func getLatestLocalChange(repoUrl string) (time.Time, error) {
	repoPath := getFolderPath(repoUrl)
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return time.Time{}, err
	}

	ref, err := repo.Head()
	if err != nil {
		return time.Time{}, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return time.Time{}, err
	}
	fmt.Println("Last local change was", commit.Author.When)
	return commit.Author.When, nil

}
*/

func NewPullService(l *logr.Logger, config *appConfig.AppConfig, c *clients.PullConfigClient, s *AppService) *PullService {
	return &PullService{
		configClient: c,
		log:          l,
		appConfig:    config,
		appService:   s,
	}
}
