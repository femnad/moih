package cmd

type KeyCfg struct {
	BucketName        string `arg:"required,-b" help:"The bucket to use"`
	CredentialFile    string `arg:"-c" help:"GCP credentials file"`
	KeySecret         string `arg:"required,-p" help:"a pass secret storing a symmetric key"`
	PasswordManager   string `arg:"-P" default:"pass" help:"Password manager, pass or 1password"`
	PrivateKey        string `arg:"-f" default:"$HOME/.ssh/{{ hostname }}" help:"a private SSH key file"`
	PrivateObjectName string `arg:"-o" default:"key/{{ hostname }}/private" help:"Object name for the private key file"`
	PublicKey         string `arg:"-l" default:"$HOME/.ssh/{{ hostname }}.pub" help:"a private SSH key file"`
	PublicObjectName  string `arg:"-u" default:"key/{{ hostname }}/public" help:"Object name for the public key file"`
}

type UpdateCfg struct {
	ApiTokenSecret  string `arg:"required,-a" help:"Git(Hub|Lab) pass secret containing API token with admin:public_key permissions"`
	KeyFile         string `arg:"-f" default:"$HOME/.ssh/{{ hostname }}.pub" help:"the public key file to upload"`
	KeyName         string `arg:"-n" default:"{{ hostname }}" help:"Key name as listed in Git(Hub|Lab)"`
	PasswordManager string `arg:"-p" default:"pass" help:"Password manager, pass or 1password"`
	Target          string `arg:"required,-t" help:"target, gitlab or github"`
	User            string `arg:"-u" default:"{{ username }}" help:"Git(Hub|Lab) username"`
}
