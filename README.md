# moih

Don't worry about deleted local SSH privates anymore, because no worries!

You can always symmetrically encrypt private keys, upload them to a GCP storage bucket and download them when required. Also update public keys on GitHub and GitLab, with even less mouse usage!

## Usage

### Download a Private Key

```
moih get [--credentialfile CREDENTIALFILE] [--objectname OBJECTNAME] --bucketname BUCKETNAME --keysecret KEYSECRET [--privatekey PRIVATEKEY]
```

### Upload a Private Key

```
moih put [--credentialfile CREDENTIALFILE] [--objectname OBJECTNAME] --bucketname BUCKETNAME --keysecret KEYSECRET [--privatekey PRIVATEKEY]
```

#### Common Options for `get` and `put`

* `--credentialfile CREDENTIALFILE`, `-c CREDENTIALFILE`
                         GCP credentials file
* `--objectname OBJECTNAME`, `-o OBJECTNAME`
                         Object name for the key file [default: private/{{ hostname }}]
* `--bucketname BUCKETNAME`, `-b BUCKETNAME`
                         The bucket to use
* `--keysecret KEYSECRET`, `-p KEYSECRET`
                         a pass[^1] secret storing a symmetric key
* `--privatekey PRIVATEKEY`, `-f PRIVATEKEY`
                         a private SSH key file [default: $HOME/.ssh/{{ hostname }}]

### Create or Update a Key on GitHub

```
moih update --apitokensecret APITOKENSECRET [--keyfile KEYFILE] [--keyname KEYNAME] [--user USER] --target TARGET
```

* `--apitokensecret APITOKENSECRET`, `-a APITOKENSECRET`
                         Git(Hub|Lab) pass secret containing API token with public_key permissions
* `--keyfile KEYFILE`, `-f KEYFILE`
                         the public key file to upload [default: private/{{ hostname }}]
* `--keyname KEYNAME`, `-n KEYNAME`
                         Key name as listed in Git(Hub|Lab) [default: {{ hostname }}]
* `--user USER`, `-u USER`
                         GitHub username [default: {{ username }}]
* `--target TARGET`, `-t TARGET`
                         target, gitlab or github

[^1]: https://www.passwordstore.org/
