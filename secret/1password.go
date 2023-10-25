package secret

type OnePassword struct {
}

func (o OnePassword) ReadSecret(secret string) (string, error) {
	return "", nil
}
