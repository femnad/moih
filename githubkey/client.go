package githubkey

import (
	"context"

	"github.com/google/go-github/v29/github"
)

func UpdateKey(token, username, title, keyContent string) (err error) {
	transport := github.BasicAuthTransport{
		Username: username,
		Password: token,
	}

	client := github.NewClient(transport.Client())
	users := client.Users

	key := github.Key{
		Key:   &keyContent,
		Title: &title,
	}
	ctx := context.Background()

	_, _, err = users.CreateKey(ctx, &key)
	if err != nil {
		return err
	}

	return
}
