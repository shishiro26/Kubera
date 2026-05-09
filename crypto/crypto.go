package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/argon2"
)

const (
	memory      uint32 = 64 * 1024
	iterations  uint32 = 3
	parallelism uint8  = 4
	keyLength   uint32 = 32
	saltLength  uint32 = 32
)

func DeriveKey(password string, salt []byte) ([]byte, error) {
	return argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength), nil
}

func Encrypt(plainText []byte, password string) (saltHex string, cipherText []byte, err error) {
	salt := make([]byte, saltLength)
	if _, err = rand.Read(salt); err != nil {
		return
	}

	key, _ := DeriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return
	}

	cipherText = gcm.Seal(nonce, nonce, plainText, nil)
	saltHex = hex.EncodeToString(salt)
	return
}

func Decrypt(saltHex string, ciphertext []byte, password string) ([]byte, error) {
	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return nil, errors.New("invalid vault format")
	}

	key, err := DeriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("corrupted vault: ciphertext too short")
	}

	nonce, data := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, errors.New("wrong master password or corrupted vault")
	}

	return plaintext, nil
}
