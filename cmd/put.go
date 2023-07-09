package cmd

import (
	"github.com/femnad/moih/gcpstorage"
	"github.com/femnad/moih/symmetric"
)

func put(cfg KeyCfg, objectName, fileName string) error {
	key, err := getSecretKey(cfg.KeySecret)
	if err != nil {
		return err
	}

	keyFile, err := expandTemplate(fileName)
	if err != nil {
		return err
	}

	encrypted, err := symmetric.Encrypt(key, keyFile)
	if err != nil {
		return err
	}

	objName, err := expandTemplate(objectName)
	if err != nil {
		return err
	}

	storageAsset := gcpstorage.StorageAsset{
		BucketName:      cfg.BucketName,
		ObjectName:      objName,
		CredentialsFile: cfg.CredentialFile,
	}

	return gcpstorage.Upload(storageAsset, encrypted)
}

func Put(cfg KeyCfg) error {
	err := put(cfg, cfg.PrivateObjectName, cfg.PrivateKey)
	if err != nil {
		return err
	}

	return put(cfg, cfg.PublicObjectName, cfg.PublicKey)
}
