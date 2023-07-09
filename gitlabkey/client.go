package gitlabkey

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"io"
)

func UpdateKey(token, title, keyContent string) error {
	git, err := gitlab.NewClient(token)
	if err != nil {
		return err
	}
	opts := gitlab.AddSSHKeyOptions{Title: &title, Key: &keyContent}
	_, resp, err := git.Users.AddSSHKey(&opts)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading body with response %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("error adding SSH key with response %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
