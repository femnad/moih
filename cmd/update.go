package cmd

import (
	"fmt"
	"os"

	"github.com/femnad/moih/githubkey"
	"github.com/femnad/moih/gitlabkey"
)

const (
	gitHub = "github"
	gitLab = "gitlab"
)

func updateGitHub(cfg UpdateCfg, key string, apiSecret string) error {
	username, err := expandTemplate(cfg.User)
	if err != nil {
		return err
	}

	title, err := expandTemplate(cfg.KeyName)
	if err != nil {
		return err
	}

	return githubkey.UpdateKey(apiSecret, username, title, key)
}

func updateGitLab(cfg UpdateCfg, key string, apiSecret string) error {
	title, err := expandTemplate(cfg.KeyName)
	if err != nil {
		return err
	}

	return gitlabkey.UpdateKey(apiSecret, title, key)
}

func Update(cfg UpdateCfg) error {
	keyFile, err := expandTemplate(cfg.KeyFile)
	if err != nil {
		return fmt.Errorf("error expanding key file name template: %v", err)
	}

	pubKey, err := os.ReadFile(keyFile)
	if err != nil {
		return fmt.Errorf("error reading public key file %s: %v", keyFile, err)
	}
	keyContent := string(pubKey)

	secretKey, err := getSecretKey(cfg.ApiTokenSecret)
	if err != nil {
		return fmt.Errorf("error looking up secret %s: %v", cfg.ApiTokenSecret, err)
	}
	apiSecret := string(secretKey)

	switch cfg.Target {
	case gitHub:
		return updateGitHub(cfg, keyContent, apiSecret)
	case gitLab:
		return updateGitLab(cfg, keyContent, apiSecret)
	default:
		return fmt.Errorf("unknown target: %s", cfg.Target)
	}
}
