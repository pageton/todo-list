package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/zalando/go-keyring"
)

func EncryptTask(plainText, base64Key string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", err
	}

	if len(key) != 32 && len(key) != 24 && len(key) != 16 {
		return "", errors.New("invalid key size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptTask(cipherText, base64Key string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", err
	}

	if len(key) != 32 && len(key) != 24 && len(key) != 16 {
		return "", errors.New("invalid key size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("cipherText too short")
	}

	nonce, encryptedData := data[:nonceSize], data[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func GenerateKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate key: %v", err)
	}

	return base64.StdEncoding.EncodeToString(key), nil
}

func StorePrivateKey(userID, privateKey string) error {
	return keyring.Set("todo-list", userID, privateKey)
}

func RetrievePrivateKey(userID string) (string, error) {
	return keyring.Get("todo-list", userID)
}
