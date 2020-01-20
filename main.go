package main

import (
	"github.com/alexflint/go-arg"
	"github.com/femnad/moih/gcpstorage"
	"github.com/femnad/moih/githubkey"
	"github.com/femnad/moih/passread"
	"github.com/femnad/moih/symmetric"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const (
	version = "0.1.0"
)

type BucketTarget struct {
	CredentialFile string `arg:"-c"`
	ObjectName string `arg:"required,-k"`
	BucketName string `arg:"required,-b"`
}

type EncryptionTarget struct {
	KeySecret  string `arg:"required,-p" help:"a pass secret storing a key"`
	SecretFile string `arg:"required,-f" help:"a file containing a private SSH key"`
}

type GetCmd struct {
	BucketTarget
	EncryptionTarget
}

type PutCmd struct {
	BucketTarget
	EncryptionTarget
}

type UpdateCmd struct {
	ApiToken string `arg:"env:API_TOKEN,required"`
	KeyFile string `arg:"required,-f"`
	KeyName string `arg:"required,-n"`
	User string `arg:"env:USER,required"`
}

type Base struct {}

func (Base) Version() string {
	return version
}

var args struct {
	Base
	Get *GetCmd `arg:"subcommand:get" help:"get a key"`
	Put *PutCmd `arg:"subcommand:put" help:"put a key"`
	Update *UpdateCmd `arg:"subcommand:update" help:"update a key in GitHub"`
}

func mustSucceed(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getSecretKey(keySecret string) []byte {
	key, err := passread.GetPassSecret(keySecret)
	mustSucceed(err)
	return []byte(key)
}

func main()  {
	p := arg.MustParse(&args)
	switch {
	case args.Put != nil:
		put := args.Put
		key := getSecretKey(put.KeySecret)
		encrypted, err := symmetric.Encrypt(key, put.SecretFile)
		mustSucceed(err)
		storageAsset := gcpstorage.StorageAsset{
			BucketName:      put.BucketName,
			ObjectName:      put.ObjectName,
			CredentialsFile: put.CredentialFile,
		}
		err = gcpstorage.Upload(storageAsset, encrypted)
		mustSucceed(err)
	case args.Get != nil:
		get := args.Get
		storageAsset := gcpstorage.StorageAsset{
			BucketName:      get.BucketName,
			ObjectName:      get.ObjectName,
			CredentialsFile: get.CredentialFile,
		}
		content, err := gcpstorage.Download(storageAsset)
		key := getSecretKey(get.KeySecret)
		decrypted, err := symmetric.Decrypt(key, content)
		mustSucceed(err)

		outputParent := path.Dir(get.SecretFile)
		if _, err = os.Stat(outputParent); os.IsNotExist(err) {
			err := os.MkdirAll(outputParent, 0700)
			mustSucceed(err)
		}
		err = ioutil.WriteFile(get.SecretFile, decrypted, 0600)
		mustSucceed(err)
	case args.Update != nil:
		update := args.Update
		err := githubkey.UpdateKey(update.ApiToken, update.User, update.KeyName, update.KeyFile)
		mustSucceed(err)
	case true:
		p.WriteHelp(os.Stdout)
	}
}

