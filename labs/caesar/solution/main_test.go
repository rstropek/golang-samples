package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaesarCipher(t *testing.T) {
	cipher := NewCaesarCipher(3)
	encrypted, err := cipher.Encrypt("Hello World!")
	assert.NoError(t, err)
	assert.Equal(t, "Khoor Zruog!", encrypted)

	decrypted, err := cipher.Decrypt(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", decrypted)
}

func TestCaesarCipherNonAlphabeticUnchanged(t *testing.T) {
	cipher := NewCaesarCipher(3)
	encrypted, err := cipher.Encrypt("Hello, World! 123")
	assert.NoError(t, err)
	assert.Equal(t, "Khoor, Zruog! 123", encrypted)

	decrypted, err := cipher.Decrypt(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World! 123", decrypted)
}

func TestCaesarCipherShiftNormalization(t *testing.T) {
	c3 := NewCaesarCipher(3)
	c29 := NewCaesarCipher(29) // 29 == 3 (mod 26)

	e1, err := c3.Encrypt("abc XYZ")
	assert.NoError(t, err)
	e2, err := c29.Encrypt("abc XYZ")
	assert.NoError(t, err)
	assert.Equal(t, e1, e2)
}

func TestCaesarCipherErrorOnDisallowedChars(t *testing.T) {
	cipher := NewCaesarCipher(3)
	_, err := cipher.Encrypt("Hello 👋")
	assert.Error(t, err)
}
