package cmd

import (
	"github.com/femnad/moih/passread"
	"os"
)

func getPublicKey(keyFile string) (string, error) {
	content, err := os.ReadFile(keyFile)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func getSecretKey(keySecret string) ([]byte, error) {
	key, err := passread.GetPassSecret(keySecret)
	if err != nil {
		return []byte{}, err
	}

	return []byte(key), nil
}
