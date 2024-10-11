package helpers

import (
	"errors"
	"os"
	"path"
)

type deployType int

const (
	UnDeployable deployType = iota
	DeploySh
)

func DeployCheck(uuid, projectsDir string) (deployType, error) {
	pdir := path.Join(projectsDir, uuid)
	bcksDeploy := path.Join(pdir, "bcks-deploy.sh")

	_, err := os.Open(bcksDeploy)
	if errors.Is(err, os.ErrNotExist) {
		return UnDeployable, nil
	} else if err != nil {
		return UnDeployable, err
	}

	return DeploySh, nil
}
