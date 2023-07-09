package cmd

import (
	"github.com/femnad/moih/gcpstorage"
	"github.com/femnad/moih/symmetric"
	"os"
	"path"
)

func getKey(cfg KeyCfg, objectName, fileName string) error {
	objName, err := expandTemplate(objectName)
	if err != nil {
		return err
	}

	storageAsset := gcpstorage.StorageAsset{
		BucketName:      cfg.BucketName,
		ObjectName:      objName,
		CredentialsFile: cfg.CredentialFile,
	}
	content, err := gcpstorage.Download(storageAsset)
	if err != nil {
		return err
	}

	key, err := getSecretKey(cfg.KeySecret)
	if err != nil {
		return err
	}

	decrypted, err := symmetric.Decrypt(key, content)
	if err != nil {
		return err
	}

	keyFile, err := expandTemplate(fileName)
	if err != nil {
		return err
	}

	outputParent := path.Dir(keyFile)
	if _, err = os.Stat(outputParent); os.IsNotExist(err) {
		mErr := os.MkdirAll(outputParent, 0700)
		if mErr != nil {
			return mErr
		}
	}

	return os.WriteFile(keyFile, decrypted, 0600)
}

func Get(cfg KeyCfg) error {
	err := getKey(cfg, cfg.PrivateObjectName, cfg.PrivateKey)
	if err != nil {
		return err
	}

	return getKey(cfg, cfg.PublicObjectName, cfg.PublicKey)
}
