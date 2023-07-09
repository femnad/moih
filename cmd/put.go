package cmd

import (
	"github.com/femnad/moih/gcpstorage"
	"github.com/femnad/moih/symmetric"
)

func Put(info PrivateKeyInfo) error {
	key, err := getSecretKey(info.KeySecret)
	if err != nil {
		return err
	}

	privKey, err := expandTemplate(info.PrivateKey)
	if err != nil {
		return err
	}

	encrypted, err := symmetric.Encrypt(key, privKey)
	if err != nil {
		return err
	}

	objName, err := expandTemplate(info.ObjectName)
	if err != nil {
		return err
	}

	storageAsset := gcpstorage.StorageAsset{
		BucketName:      info.BucketName,
		ObjectName:      objName,
		CredentialsFile: info.CredentialFile,
	}

	return gcpstorage.Upload(storageAsset, encrypted)
}
