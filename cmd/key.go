package cmd

import (
	"fmt"
	"github.com/femnad/moih/secret"
)

func getSecretKey(pwMgr, keySecret string) ([]byte, error) {
	mgr, err := passwordManager(pwMgr)
	if err != nil {
		return []byte{}, err
	}

	key, err := mgr.ReadSecret(keySecret)
	if err != nil {
		return []byte{}, err
	}

	return []byte(key), nil
}

func passwordManager(pwMgr string) (secret.Manager, error) {
	switch pwMgr {
	case "pass":
		return secret.Pass{}, nil
	case "1password":
		return &secret.OnePassword{}, nil
	default:
		return nil, fmt.Errorf("unknown password manager: %s", pwMgr)
	}
}
