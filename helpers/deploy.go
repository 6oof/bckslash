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
	DeployPlain
)

func DeployCheck(uuid string) (deployType, error) {
	pdir := path.Join("projects", uuid)
	bcksCompose := path.Join(pdir, "bckslash-compose.yaml")

	_, err := os.Open(bcksCompose)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return UnDeployable, errors.New("bckslash-compose.yaml does not exist.\nmake sure you follow the getting started guide")
		}
		return UnDeployable, err
	}

	bcksDeploy := path.Join(pdir, "bcks-deploy.sh")
	hasDeploySh := false

	_, err = os.Open(bcksDeploy)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return UnDeployable, err
	} else {
		hasDeploySh = true
	}

	if hasDeploySh {
		return DeploySh, nil

	} else {
		return DeployPlain, nil
	}

}
