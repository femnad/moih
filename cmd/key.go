package cmd

import (
	"github.com/femnad/moih/secret"
)

func getSecretKey(keySecret string) ([]byte, error) {
	mgr := passwordManager()
	key, err := mgr.ReadSecret(keySecret)
	if err != nil {
		return []byte{}, err
	}

	return []byte(key), nil
}

func passwordManager() secret.Manager {
	return secret.Pass{}
}
