package githubkey

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

type sshKey struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// For some reason keys returned from github.Client.Users.ListKeys don't have title set.
func lookupUserKeys(token string) ([]sshKey, error) {
	var keys []sshKey

	hc := http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user/keys", nil)
	if err != nil {
		return keys, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	resp, err := hc.Do(req)
	if err != nil {
		return keys, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return keys, fmt.Errorf("error reading body of response with status %d: %v", resp.StatusCode, err)
		}

		return keys, fmt.Errorf("unexpected response with status %d and body %s", resp.StatusCode, body)
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&keys)
	if err != nil {
		return keys, err
	}

	return keys, nil
}

func UpdateKey(token, username, title, keyContent string) (err error) {
	ctx := context.TODO()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	users := client.Users

	keys, err := lookupUserKeys(token)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if key.Title != title {
			continue
		}

		_, err = users.DeleteKey(ctx, key.ID)
		if err != nil {
			return err
		}
		break
	}

	key := github.Key{
		Key:   &keyContent,
		Title: &title,
	}

	_, _, err = users.CreateKey(ctx, &key)
	if err != nil {
		return err
	}

	return
}
