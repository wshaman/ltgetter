package tools

import (
	"errors"
	"os"
)

const ltEnv = "LITNET_IDENTITY"

func GetIdentity() (identity string, err error) {
	if identity = os.Getenv(ltEnv); identity == "" {
		return "", errors.New("no identity set ")
	}
	return identity, nil
}
