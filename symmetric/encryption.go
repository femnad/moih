package symmetric

// Shameless copy paste from https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
)

func getCipher(key []byte) (cipher.AEAD, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error initializing cipher: %s", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("error using Galois/Counter mode: %s", err)
	}

	return gcm, nil
}

func Encrypt(key []byte, secretFileName string) (output []byte, err error) {
	content, err := ioutil.ReadFile(secretFileName)
	if err != nil {
		return output, fmt.Errorf("error read file contents from %s: %s", secretFileName, err)
	}

	c, err := getCipher(key)
	if err != nil {
		return output, fmt.Errorf("error getting a cipher: %s", err)
	}

	nonce := make([]byte, c.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return output, fmt.Errorf("error populating nonce with random values: %s", err)
	}

	output = c.Seal(nonce, nonce, content, nil)
	return
}

func Decrypt(key []byte, encrypted []byte) (output []byte, err error){
	c, err := getCipher(key)
	if err != nil {
		return output, fmt.Errorf("error getting a cipher: %s", err)
	}

	c, err = getCipher(key)
	if err != nil {
		return output, fmt.Errorf("error getting a cipher: %s", err)
	}

	nonceSize := c.NonceSize()
	if len(encrypted) < nonceSize {
		return output, fmt.Errorf("size of ciphertext (%d) smaller than nonce size (%d)", len(encrypted), nonceSize)
	}

	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]
	output, err = c.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return output, fmt.Errorf("error decrypting ciphertext: %s", err)
	}
	return
}
