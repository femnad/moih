# moih

Don't worry about deleted local SSH privates anymore, because no worries!

Or you can always symettrically encrypt private key, upload them to a GCP storage bucket and download them when required. Also update public keys on GitHub, with even less mouse usage!

## Usage

### Download a Private Key

```
moih get [-c CREDENTIALFILE] -o OBJECTNAME -b BUCKETNAME -p KEYSECRET -p SECRETFILE
```

### Upload a Private Key

```
moih put [-c CREDENTIALFILE] -o OBJECTNAME -b BUCKETNAME -p KEYSECRET -p SECRETFILE
```

#### Common Options for `get` and `put`

* `KEYSECRET`: A [pass](pass) secret storing a symmetric key
* `CREDENTIALFILE`: A Google Cloud credentials file, if not provided default to the usual credential lookup

### Create or Update a Key on GitHub

```
moih update -f KEYFILE -n KEYNAME
```

also as environment variables

* `API_TOKEN`: GitHub API token with `admin:public_key` permissions
* `USER`: GitHub user, can be overridden not to use the OS user with `-u`

[pass]: https://www.passwordstore.org/
