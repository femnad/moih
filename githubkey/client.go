package githubkey

import (
	"context"
	"fmt"
	"github.com/google/go-github/v29/github"
	"io/ioutil"
)

func UpdateKey(token, username, title, keyFile string) (err error) {
	transport := github.BasicAuthTransport{
		Username:  username,
		Password:  token,
	}

	client := github.NewClient(transport.Client())
	users := client.Users

	content, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return fmt.Errorf("error reading file %s: %s", keyFile, err)
	}

	publicKey := string(content)

	key := github.Key{
		Key:       &publicKey,
		Title:     &title,
	}
	ctx := context.Background()

	_, _, err = users.CreateKey(ctx, &key)
	if err != nil {
		return err
	}

	return
}

