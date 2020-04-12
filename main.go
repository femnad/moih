package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/alexflint/go-arg"
	"github.com/femnad/moih/gcpstorage"
	"github.com/femnad/moih/githubkey"
	"github.com/femnad/moih/gitlabkey"
	"github.com/femnad/moih/passread"
	"github.com/femnad/moih/symmetric"
)

const (
	version = "0.1.3"
	GitHub  = "github"
	GitLab  = "gitlab"
)

type BucketTarget struct {
	CredentialFile string `arg:"-c"`
	ObjectName     string `arg:"required,-o"`
	BucketName     string `arg:"required,-b"`
}

type EncryptionTarget struct {
	KeySecret  string `arg:"required,-p" help:"a pass secret storing a symmetric key"`
	SecretFile string `arg:"required,-f" help:"a private SSH key file"`
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
	ApiToken string `arg:"env:API_TOKEN,required,-a" help:"Git(Hub|Lab) API token with admin:public_key permissions"`
	KeyFile  string `arg:"required,-f" help:"the public key file to upload"`
	KeyName  string `arg:"required,-n" help:"Key name as listed in Git(Hub|Lab)"`
	User     string `arg:"env:USER,required,-u" help:"GitHub username"`
	Target   string `arg:"required,-t" help:"target, gitlab or github"`
}

func (u UpdateCmd) updateGitHub(key string) {
	err := githubkey.UpdateKey(u.ApiToken, u.User, u.KeyName, key)
	mustSucceed(err)
}

func (u UpdateCmd) updateGitLab(key string) {
	err := gitlabkey.UpdateKey(u.ApiToken, u.KeyName, key)
	mustSucceed(err)
}

type Base struct{}

func (Base) Version() string {
	return version
}

var args struct {
	Base
	Get    *GetCmd    `arg:"subcommand:get" help:"get a key from GCP Cloud Storage"`
	Put    *PutCmd    `arg:"subcommand:put" help:"put a key into GCP Cloud Storage"`
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

func getPublicKey(keyFile string) string {
	content, err := ioutil.ReadFile(keyFile)
	mustSucceed(err)
	return string(content)
}

func main() {
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
		mustSucceed(err)
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
		key := getPublicKey(update.KeyFile)
		switch {
		case update.Target == GitHub:
			update.updateGitHub(key)
		case update.Target == GitLab:
			update.updateGitLab(key)
		case update.Target == "":
			update.updateGitHub(key)
			update.updateGitLab(key)
		default:
			log.Fatalf("Unknown target: %s", update.Target)
		}
	case true:
		p.WriteHelp(os.Stdout)
	}
}
