package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

type AES struct {
	Key []byte
}

func NewAES(key string) (*AES, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length")
	}
	
	return &AES{Key: []byte(key)}, nil
}

func (a *AES) Encrypt(data string) (encryptedText []byte, err error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return encryptedText, err
	}

	cipherText := aesGCM.Seal(nil, nonce, []byte(data), nil)
	encryptedText = append(nonce, cipherText...)
	return
}

func (a *AES) Decrypt(data string) (decryptedText []byte, err error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	encryptedBytes := []byte(data)
	nonce := encryptedBytes[:aesGCM.NonceSize()]
	cipherText := encryptedBytes[aesGCM.NonceSize():]

	decryptedText, err = aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return
	}
	return
}
