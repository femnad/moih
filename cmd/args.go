package cmd

type PrivateKeyInfo struct {
	CredentialFile string `arg:"-c" help:"GCP credentials file"`
	ObjectName     string `arg:"-o" default:"private/{{ hostname }}" help:"Object name for the key file"`
	BucketName     string `arg:"required,-b" help:"The bucket to use"`
	KeySecret      string `arg:"required,-p" help:"a pass secret storing a symmetric key"`
	PrivateKey     string `arg:"-f" default:"$HOME/.ssh/{{ hostname }}" help:"a private SSH key file"`
}

type UpdateCfg struct {
	ApiTokenSecret string `arg:"required,-a" help:"Git(Hub|Lab) pass secret containing API token with admin:public_key permissions"`
	KeyFile        string `arg:"-f" help:"the public key file to upload"`
	KeyName        string `arg:"-n" default:"{{ hostname }}" help:"Key name as listed in Git(Hub|Lab)"`
	User           string `arg:"-u" default:"{{ username }}" help:"Git(Hub|Lab) username"`
	Target         string `arg:"required,-t" help:"target, gitlab or github"`
}
