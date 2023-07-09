package cmd

import (
	"fmt"
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
		return err
	}

	pubKey, err := getPublicKey(keyFile)
	if err != nil {
		return err
	}

	secretKey, err := getSecretKey(cfg.ApiTokenSecret)
	if err != nil {
		return err
	}

	apiSecret := string(secretKey)
	switch cfg.Target {
	case gitHub:
		return updateGitHub(cfg, pubKey, apiSecret)
	case gitLab:
		return updateGitLab(cfg, pubKey, apiSecret)
	default:
		return fmt.Errorf("unknown target: %s", cfg.Target)
	}
}
