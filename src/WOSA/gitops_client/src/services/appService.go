package services

import (
	"errors"
	"fmt"
	"gitops_client/appConfig"
	"gitops_client/clients"
	"gitops_client/models"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-logr/logr"
)

const (
	solutionFileName = "solution.yaml"
	instanceFileName = "instance.yaml"
)

type AppService struct {
	log           *logr.Logger
	appConfig     *appConfig.AppConfig
	kubectlClient *clients.KubectlClient
}

func (s AppService) ProcessRepo() {
	d, err := s.GetDesiredState()
	if err != nil {
		s.log.Error(err, "Error trying to get the desired state")
	}

	if err := s.ReconcileDesirdState(d); err != nil {
		s.log.Error(err, "Error reconciling desired state")
	}
}

func (s AppService) ReconcileDesirdState(d []models.DesiredState) error {
	s.log.Info("Reconciling desired state...")
	s.log.Info("TODO: Need to handle multiple apps")

	solutionIdx := slices.IndexFunc(d, func(s models.DesiredState) bool { return s.FileName == solutionFileName })
	instanceIdx := slices.IndexFunc(d, func(s models.DesiredState) bool { return s.FileName == instanceFileName })
	if solutionIdx < 0 || instanceIdx < 0 {
		s.log.Info("Solution and Instance files not found. Processing stopped")
		return nil
	}

	desiredState := d[solutionIdx]
	if err := s.ApplyAndCopy(desiredState, models.KindSolution); err != nil {
		return nil
	}

	desiredState = d[instanceIdx]
	if err := s.ApplyAndCopy(desiredState, models.KindInstance); err != nil {
		return nil
	}

	for i, v := range d {
		if v.State == models.Unchanged {
			continue
		}

		if i == solutionIdx || i == instanceIdx {
			continue
		}

		if err := s.CopyFile(v); err != nil {
			return err
		}
	}

	return nil
}

func (s AppService) ApplyAndCopy(d models.DesiredState, kind string) error {
	s.log.Info(fmt.Sprintf("Applying and copying desired state %+v", d))

	if d.State != models.Unchanged {
		err := s.kubectlClient.Apply(d, kind)
		if err != nil {
			s.log.Error(err, fmt.Sprintf("Error applying desired state for file %s/%s", d.CurrentStatePath, d.FileName))
			return err
		}
		s.CopyFile(d)
	}
	return nil
}

func (s AppService) CopyFile(d models.DesiredState) error {
	s.log.Info(fmt.Sprintf("Copying desired state %+v", d))

	if d.State == models.Unchanged {
		return nil
	}

	var dst *os.File
	var err error
	dstPath := path.Join(d.CurrentStatePath, d.FileName)

	_, err = os.Stat(d.CurrentStatePath)
	if err != nil {
		s.log.Info("Creating current state app directory")
		if err := os.MkdirAll(d.CurrentStatePath, os.ModePerm); err != nil {
			s.log.Error(err, "Error creating current state app directory")
			return err
		}
	}

	if d.State == models.New {
		dst, err = os.Create(dstPath)
		if err != nil {
			s.log.Error(err, fmt.Sprintf("Error creating new current state file %s", dstPath))
			return err
		}
	} else {
		dst, err = os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			s.log.Error(err, fmt.Sprintf("Error creating opening current state file %s", dstPath))
			return err
		}
	}
	defer dst.Close()

	src, err := os.Open(d.SourceFile)
	if err != nil {
		s.log.Error(err, "Error opening desired state file %s", d.SourceFile)
		return err
	}
	defer src.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		s.log.Error(err, fmt.Sprintf("Error copying desired state file %s to current state location %s", d.SourceFile, dstPath))
		return err
	}

	err = dst.Sync()
	if err != nil {
		s.log.Error(err, fmt.Sprintf("Error saving contents of the desired state file %s", dstPath))
		return err
	}

	return nil
}

func (s AppService) GetDesiredState() ([]models.DesiredState, error) {
	repoPath := s.appConfig.GetRepoFolder()
	currentStatePath := s.appConfig.GetCurrentStateFolder()
	_, err := os.Stat(currentStatePath)
	if err != nil {
		s.log.Info("Creating current state device directory")
		if err := os.MkdirAll(currentStatePath, os.ModePerm); err != nil {
			s.log.Error(err, "Error creating current state device directory")
			return nil, err
		}
	}

	desiredState := []models.DesiredState{}
	err = filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if strings.Contains(path, ".git") {
			return nil
		}

		s.log.Info(fmt.Sprintf("Checking path %s", path))

		if d.IsDir() {
			return nil
		}

		currentStateFile := strings.Replace(path, s.appConfig.DeviceRepoRootPath, s.appConfig.CurrentStateRootPath, -1)
		s.log.Info(fmt.Sprintf("Checking if file %s exists", currentStateFile))
		cfi, fErr := os.Stat(currentStateFile)
		if fErr != nil {
			if errors.Is(fErr, os.ErrNotExist) {
				desiredState = append(desiredState, models.DesiredState{
					FileName:         d.Name(),
					CurrentStatePath: filepath.Dir(currentStateFile),
					SourceFile:       path,
					State:            models.New,
				})
				return nil
			} else {
				s.log.Error(fErr, "Error checking current state folder")
				return fErr
			}
		}

		rfi, fErr := d.Info()
		if fErr != nil {
			s.log.Error(fErr, "Error getting file info")
			return fErr
		}

		var state models.FileState
		if os.SameFile(rfi, cfi) {
			state = models.Unchanged
		} else {
			state = models.Updated
		}

		desiredState = append(desiredState, models.DesiredState{
			FileName:         d.Name(),
			CurrentStatePath: filepath.Dir(currentStateFile),
			SourceFile:       path,
			State:            state,
		})

		return nil
	})
	if err != nil {
		s.log.Error(err, "Error processing files")
		return nil, err
	}

	return desiredState, nil
}

func NewAppService(l *logr.Logger, c *appConfig.AppConfig, k *clients.KubectlClient) *AppService {
	return &AppService{
		log:           l,
		appConfig:     c,
		kubectlClient: k,
	}
}
