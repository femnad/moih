package main

import (
	"bytes"
	"html/template"
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
	version = "0.3.2"
	GitHub  = "github"
	GitLab  = "gitlab"
)

type PrivateKeyInfo struct {
	CredentialFile string `arg:"-c" help:"GCP credentials file"`
	ObjectName     string `arg:"-o" default:"private/{{ hostname }}" help:"Object name for the key file"`
	BucketName     string `arg:"required,-b" help:"The bucket to use"`
	KeySecret      string `arg:"required,-p" help:"a pass secret storing a symmetric key"`
	PrivateKey     string `arg:"-f" default:"$HOME/.ssh/{{ hostname }}" help:"a private SSH key file"`
}

type GetCmd struct {
	PrivateKeyInfo
}

type PutCmd struct {
	PrivateKeyInfo
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

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
	    log.Fatalf("error getting hostname: %v", err)
	}

	return hostname
}

func expandTemplate(text string) string {
	text = os.ExpandEnv(text)

	tmpl := template.New("moih")
	tmpl.Funcs(map[string]interface{}{"hostname": getHostname})

	parsed, err := tmpl.Parse(text)
	if err != nil {
	    log.Fatalf("error parsing template %s: %v", text, err)
	}

	out := bytes.Buffer{}
	err = parsed.Execute(&out, struct{}{})
	if err != nil {
	    log.Fatalf("error executing template %s: %v", text, err)
	}

	return out.String()
}

func main() {
	p := arg.MustParse(&args)
	switch {
	case args.Put != nil:
		put := args.Put
		key := getSecretKey(put.KeySecret)
		encrypted, err := symmetric.Encrypt(key, expandTemplate(put.PrivateKey))
		mustSucceed(err)
		storageAsset := gcpstorage.StorageAsset{
			BucketName:      put.BucketName,
			ObjectName:      expandTemplate(put.ObjectName),
			CredentialsFile: put.CredentialFile,
		}
		err = gcpstorage.Upload(storageAsset, encrypted)
		mustSucceed(err)
	case args.Get != nil:
		get := args.Get
		storageAsset := gcpstorage.StorageAsset{
			BucketName:      get.BucketName,
			ObjectName:      expandTemplate(get.ObjectName),
			CredentialsFile: get.CredentialFile,
		}
		content, err := gcpstorage.Download(storageAsset)
		mustSucceed(err)
		key := getSecretKey(get.KeySecret)
		decrypted, err := symmetric.Decrypt(key, content)
		mustSucceed(err)

		outputParent := path.Dir(expandTemplate(get.PrivateKey))
		if _, err = os.Stat(outputParent); os.IsNotExist(err) {
			err := os.MkdirAll(outputParent, 0700)
			mustSucceed(err)
		}
		err = ioutil.WriteFile(expandTemplate(get.PrivateKey), decrypted, 0600)
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
