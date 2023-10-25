package secret

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/anmitsu/go-shlex"
)

type opFields struct {
	Value string `json:"value"`
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

	scr := opSecret{}
	decoder := json.NewDecoder(stdout)
	err = decoder.Decode(&scr)
	if err != nil {
		return "", err
	}

	if len(scr.Fields) == 0 {
		return "", fmt.Errorf("cannot find first section in secret %s", secret)
	}

	return scr.Fields[0].Value, nil
}
