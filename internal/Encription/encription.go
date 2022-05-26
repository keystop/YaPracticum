package encription

import (
	"crypto/aes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
)

var keySize int = 16

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

func EncriptStr(s string) (string, error) {
	key, err := generateRandom(aes.BlockSize)
	if err != nil {
		return "", err
	}

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	res := make([]byte, aesblock.BlockSize())
	hash := md5.Sum([]byte(s))
	aesblock.Encrypt(res, hash[:])

	return fmt.Sprintf("%x", res), nil

}
