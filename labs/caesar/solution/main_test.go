package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaesarCipher(t *testing.T) {
	cipher := NewCaesarCipher(3)
	encrypted, err := cipher.Encrypt("Hello World!")
	assert.Nil(t, err)
	assert.Equal(t, encrypted, "Khoor Zruog!")

	decrypted, err := cipher.Decrypt(encrypted)
	assert.Nil(t, err)
	assert.Equal(t, decrypted, "Hello World!")
}

func TestCeasarCipherErrorWithSpecialChars(t *testing.T) {
	cipher := NewCaesarCipher(3)
	_, err := cipher.Encrypt("Hello World 👋")
	assert.NotNil(t, err)
}
