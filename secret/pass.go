package secret

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Pass struct {
}

func (p Pass) ReadSecret(secret string) (out string, err error) {
	cmd := exec.Command("pass", secret)
	stdout := bytes.Buffer{}
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		return out, fmt.Errorf("error getting pass secret %s: %s", secret, err)
	}
	outBytes, err := io.ReadAll(&stdout)
	out = strings.TrimSpace(string(outBytes))
	return
}
