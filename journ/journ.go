package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

var (
	plaintextInputFile  = "/Users/seanmccullough/code/me3/journ/plaintext-input.txt"
	ciphertextFile      = "/Users/seanmccullough/code/me3/journ/ciphertext.bin"
	keyFile             = "/Users/seanmccullough/code/me3/journ/key.txt"
	plaintextOutputFile = "/Users/seanmccullough/code/me3/journ/plaintext-output.txt"
)

func encryptFile() error {
	plainText, err := os.ReadFile(plaintextInputFile)
	if err != nil {
		return err
	}

	key, err := os.ReadFile(keyFile)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

	if err := os.WriteFile(ciphertextFile, cipherText, 0777); err != nil {
		return err
	}

	return nil
}

func decryptFile() error {
	cipherText, err := os.ReadFile(ciphertextFile)
	if err != nil {
		return err
	}

	key, err := os.ReadFile(keyFile)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return err
	}

	if err := os.WriteFile(plaintextOutputFile, plainText, 0777); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := encryptFile(); err != nil {
		fmt.Printf("error encrpting: %+v\n", err)
		os.Exit(1)
	}
	if err := decryptFile(); err != nil {
		fmt.Printf("error decripting: %+v\n", err)
		os.Exit(1)
	}
}
