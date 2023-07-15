package gitlabkey

import (
	"fmt"
	"io"

	"github.com/xanzy/go-gitlab"
)

func UpdateKey(token, username, title, keyContent string) error {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return err
	}

	var keyId int
	keys, resp, err := client.Users.ListSSHKeysForUser(username, &gitlab.ListSSHKeysForUserOptions{})
	for _, key := range keys {
		if key.Title == title {
			keyId = key.ID
			break
		}
	}

	if keyId != 0 {
		_, err = client.Users.DeleteSSHKey(keyId)
		if err != nil {
			return err
		}
	}

	opts := gitlab.AddSSHKeyOptions{Title: &title, Key: &keyContent}
	_, resp, err = client.Users.AddSSHKey(&opts)
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
