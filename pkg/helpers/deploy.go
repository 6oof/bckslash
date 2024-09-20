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

func DeployCheck(uuid string) (deployType, error) {
	pdir := path.Join("projects", uuid)

	bcksDeploy := path.Join(pdir, "bcks-deploy.sh")

	_, err := os.Open(bcksDeploy)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return UnDeployable, err
	} else {

		return DeploySh, nil
	}

}
