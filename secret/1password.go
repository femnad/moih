package secret

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/anmitsu/go-shlex"
)

const (
	passwordFieldLabel = "password"
)

type opFields struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type opSecret struct {
	Fields []opFields `json:"fields"`
}

type OnePassword struct {
}

func (o OnePassword) ReadSecret(secret string) (string, error) {
	cmdArgs, err := shlex.Split(fmt.Sprintf("item get %s --format json", secret), true)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("op", cmdArgs...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err = cmd.Start(); err != nil {
		return "", err
	}

	var decoded opSecret
	decoder := json.NewDecoder(stdout)
	err = decoder.Decode(&decoded)
	if err != nil {
		return "", err
	}

	for _, field := range decoded.Fields {
		if field.Label == passwordFieldLabel {
			return field.Value, nil
		}
	}

	return "", fmt.Errorf("unable to password field in secret %s", secret)
}
