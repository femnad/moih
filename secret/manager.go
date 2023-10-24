package secret

type Manager interface {
	ReadSecret(secret string) (out string, err error)
}
