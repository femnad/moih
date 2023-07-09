package cmd

import (
	"github.com/femnad/moih/gcpstorage"
	"github.com/femnad/moih/symmetric"
	"os"
	"path"
)

func Get(info PrivateKeyInfo) error {
	objName, err := expandTemplate(info.ObjectName)
	if err != nil {
		return err
	}

	storageAsset := gcpstorage.StorageAsset{
		BucketName:      info.BucketName,
		ObjectName:      objName,
		CredentialsFile: info.CredentialFile,
	}
	content, err := gcpstorage.Download(storageAsset)
	if err != nil {
		return err
	}

	key, err := getSecretKey(info.KeySecret)
	if err != nil {
		return err
	}

	decrypted, err := symmetric.Decrypt(key, content)
	if err != nil {
		return err
	}

	privKey, err := expandTemplate(info.PrivateKey)
	if err != nil {
		return err
	}

	outputParent := path.Dir(privKey)
	if _, err = os.Stat(outputParent); os.IsNotExist(err) {
		mErr := os.MkdirAll(outputParent, 0700)
		if mErr != nil {
			return mErr
		}
	}

	return os.WriteFile(privKey, decrypted, 0600)
}
