package services

import (
	"encoding/json"
	"fmt"
	"gitops_pushservice/appConfig"
	"gitops_pushservice/clients"
	"gitops_pushservice/models"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-logr/logr"
	"gopkg.in/yaml.v3"
)

type PushService struct {
	log       *logr.Logger
	appConfig *appConfig.AppConfig
	mqClient  *clients.MQClient
}

const (
	defaultRemoteName = "origin"
	solutionFileName  = "solution.yaml"
	instanceFileName  = "instance.yaml"
)

func (s PushService) Start() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	deliveries, err := s.mqClient.Consume()
	if err != nil {
		s.log.Error(err, "Could not start consuming")
		return
	}

	for {
		select {
		case <-sigCh:
			fmt.Println("Received interrupt signal. Shutting down gracefully.")
			os.Exit(0)
		case delivery := <-deliveries:
			// Ack a message every 2 seconds
			var message models.StateChanged
			err := json.Unmarshal(delivery.Body, &message)
			if err != nil {
				s.log.Info(fmt.Sprintf("Message body: %s", string(delivery.Body)))
				s.log.Error(err, "error unmarshalling message")
				delivery.Nack(false, true)
			}
			s.log.Info(fmt.Sprint("Received message: ", message))
			if message.Kind == models.Install {
				s.log.Info("Installing new app")
				s.installApp(message)
			}

			if err := delivery.Ack(false); err != nil {
				fmt.Printf("Error acknowledging message: %s\n", err)
			}
			<-time.After(time.Second * 2)
		}
	}
}

func (s PushService) installApp(newApp models.StateChanged) error {
	repoPath := path.Join(s.appConfig.RepoRootPath, newApp.DeviceId.String())
	err := s.cloneIfNew(newApp.DeviceRepo.Url, newApp.DeviceRepo.Branch, repoPath)
	if err != nil {
		return nil
	}

	appPath := path.Join(repoPath, newApp.AppId.String())
	_, err = os.Stat(appPath)
	if err != nil {
		os.Mkdir(appPath, os.ModePerm)
	}

	err = s.createSolutionFile(appPath, newApp)
	if err != nil {
		s.log.Error(err, "Error creating solution file")
		return err
	}
	s.createInstanceFile(appPath, newApp)
	if err != nil {
		s.log.Error(err, "Error creating instance file")
	}

	err = s.commitNewProject(newApp.DeviceRepo.Url, repoPath)
	if err != nil {
		return err
	}

	return nil
}

func (s PushService) createSolutionFile(appPath string, newApp models.StateChanged) error {
	filePath := path.Join(appPath, solutionFileName)
	file, err := os.Create(filePath)
	if err != nil {
		s.log.Error(err, "Error creating file: %v\n", err)
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	solution, err := s.buildSolution(newApp.Sources[0], newApp.Properties)
	if err != nil {
		s.log.Error(err, "Error building solution")
		return err
	}

	if err := encoder.Encode(solution); err != nil {
		s.log.Error(err, "Error encoding YAML")
		return err
	}

	return nil
}

func (s PushService) buildSolution(src models.Source, p map[string]models.Property) (models.Solution, error) {

	envValues := make(map[string]interface{})
	envString := ""
	for _, v := range p {
		for _, t := range v.Targets {
			p := strings.Split(t.Pointer, "/")
			envString += fmt.Sprintf("  %s: %s\n", p[len(p)-1], v.Value)
		}
	}
	s.log.Info(fmt.Sprint("EnvString: ", envString))
	envData := []byte(envString)
	s.log.Info(fmt.Sprintf("EnvString Byte Len %v", len(envData)))
	var envI interface{}
	err := yaml.Unmarshal(envData, &envI)
	if err != nil {
		s.log.Error(err, "error unmarshaling properties")
		return models.Solution{}, err
	}
	envValues["env"] = envI
	s.log.Info(fmt.Sprintf("Unmarshal result: %+v", envI))
	helmProps := models.HelmPropertyConfig{
		Repo:    src.Properties.Repository,
		Version: src.Properties.Revision,
		Values:  envValues,
		Name:    src.Name,
		Wait:    src.Properties.Wait,
	}
	solution := models.Solution{
		ApiVersion: "solution.symphony/v1",
		Kind:       "Solution",
		Metadata: models.SolutionMetadata{
			Name:      fmt.Sprintf("%s-solution", src.Name),
			Namespace: src.Name,
		},
		Spec: models.SolutionSpec{
			Components: []models.ComponentSpec{
				{
					Name:       src.Name,
					Type:       "helm.v3",
					Properties: helmProps,
				},
			},
		},
	}
	return solution, nil
}

func (s PushService) createInstanceFile(appPath string, newApp models.StateChanged) error {
	filePath := path.Join(appPath, instanceFileName)
	file, err := os.Create(filePath)
	if err != nil {
		s.log.Error(err, "Error creating file: %v\n", err)
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	instance, err := s.buildInstance(newApp.Sources[0])
	if err != nil {
		s.log.Error(err, "Error building instance")
		return err
	}

	if err := encoder.Encode(instance); err != nil {
		s.log.Error(err, "Error encoding YAML")
		return err
	}

	return nil
}

func (s PushService) buildInstance(h models.Source) (models.Instance, error) {
	var instance = models.Instance{
		ApiVersion: "solution.symphony/v1",
		Kind:       "Instance",
		Metadata: models.InstanceMetadata{
			Name:      fmt.Sprintf("%s-instance", h.Name),
			Namespace: h.Name,
		},
		Spec: models.InstanceSpec{
			Name:     h.Name,
			Scope:    h.Name,
			Solution: fmt.Sprintf("%s-solution", h.Name),
			Target: models.TargetSelector{
				Name: "unknow-target",
			},
		},
	}

	return instance, nil
}

func (s PushService) commitNewProject(repoUrl, repoPath string) error {
	s.log.Info("Committing new project")
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		s.log.Error(err, fmt.Sprint("Error opening git repo ", repoPath))
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		s.log.Error(err, "Error geting worktree")
		return err
	}

	err = w.AddWithOptions(&git.AddOptions{
		All: true,
	})
	if err != nil {
		s.log.Error(err, "Error adding solution file")
		return err
	}
	_, err = w.Commit("adding new app solution", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Push Service",
			Email: "pushservice@email.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		s.log.Error(err, "Error committing chages")
		return err
	}

	return s.pushChanges(r, repoUrl)
}

func (s PushService) cloneIfNew(repoUrl, branch, repoPath string) error {
	s.log.Info(fmt.Sprintf("Checking if device repo %s is new", repoUrl))
	_, err := os.Stat(repoPath)
	if err != nil {
		s.log.Info(fmt.Sprintf("Cloning new repo %s to %s", repoUrl, repoUrl))
		repo, err := git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:               repoUrl,
			Progress:          os.Stdout,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			s.log.Error(err, "Cloning repository")
			return err
		}

		err = s.checkoutBranch(repo, branch)
		return err
	}

	return nil
}

func (s PushService) checkoutBranch(repo *git.Repository, branch string) error {
	s.log.Info(fmt.Sprint("Checking out branch ", branch))
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

func (s PushService) pushChanges(r *git.Repository, repoUrl string) error {

	_, err := r.Remote(defaultRemoteName)
	if err != nil {
		s.log.Info(fmt.Sprint("Creating new Git remote named " + defaultRemoteName))
		_, err = r.CreateRemote(&config.RemoteConfig{
			Name: defaultRemoteName,
			URLs: []string{repoUrl},
		})
		if err != nil {
			s.log.Error(err, fmt.Sprint("Error creating remote:", err))
		}
	}

	err = r.Push(&git.PushOptions{
		RemoteName: defaultRemoteName,
		Auth: &http.BasicAuth{
			Username: s.appConfig.DeviceRepoUserName,
			Password: s.appConfig.DeviceRepoPassword,
		},
	})
	if err != nil {
		s.log.Error(err, "error pushing changes")
		return err
	}

	return nil
}

func NewPushService(l *logr.Logger, c *appConfig.AppConfig, mcl *clients.MQClient) *PushService {
	return &PushService{
		log:       l,
		appConfig: c,
		mqClient:  mcl,
	}
}
