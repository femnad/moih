package passread

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func GetPassSecret(secret string) (out string, err error) {
	cmd := exec.Command("pass", secret)
	stdout := bytes.Buffer{}
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		return out, fmt.Errorf("error getting pass secret %s: %s", secret, err)
	}
	outBytes, err := ioutil.ReadAll(&stdout)
	out = strings.TrimSpace(string(outBytes))
	return
}
