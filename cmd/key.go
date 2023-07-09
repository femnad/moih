package cmd

import (
	"github.com/femnad/moih/passread"
)

func getSecretKey(keySecret string) ([]byte, error) {
	key, err := passread.GetPassSecret(keySecret)
	if err != nil {
		return []byte{}, err
	}

	return []byte(key), nil
}
