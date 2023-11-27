package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Encryptor interface {
	Encrypt(plain string) (string, error)
}

type Decryptor interface {
	Decrypt(cypher string) (string, error)
}

type EncryptDecryptor interface {
	Encryptor
	Decryptor
}

type CaesarCipher struct {
	shift int
}

func NewCaesarCipher(shift int) EncryptDecryptor {
	return &CaesarCipher{shift: shift}
}

func (cc *CaesarCipher) Encrypt(plain string) (string, error) {
	return cc.shiftText(plain, cc.shift)
}

func (cc *CaesarCipher) Decrypt(cipher string) (string, error) {
	return cc.shiftText(cipher, -cc.shift)
}

func (cc *CaesarCipher) shiftText(text string, shift int) (string, error) {
	allowedChars := map[rune]bool{
		' ':  true,
		'.':  true,
		',':  true,
		'!':  true,
		'?':  true,
		'\'': true,
		'-'	: true,
		'\n'	: true,
	}

	var resultBuilder strings.Builder
	resultBuilder.Grow(len(text))
	for _, char := range text {
		if char >= 'a' && char <= 'z' {
			resultBuilder.WriteRune(((char-'a'+rune(shift))%26+26)%26 + 'a')
		} else if char >= 'A' && char <= 'Z' {
			resultBuilder.WriteRune(((char-'A'+rune(shift))%26+26)%26 + 'A')
		} else if allowedChars[char] {
			resultBuilder.WriteRune(char)
		} else {
			return "", fmt.Errorf("invalid character: %c", char)
		}
	}
	return resultBuilder.String(), nil
}

func main() {
	var shift int
	var mode, text string

	flag.IntVar(&shift, "shift", 0, "Shift value for the Caesar cipher")
	flag.StringVar(&mode, "mode", "encrypt", "Mode of operation: encrypt or decrypt")
	flag.StringVar(&text, "text", "", "Text to process")
	flag.Parse()

	cipher := NewCaesarCipher(shift)

	// If text starts with @, read from file
	if len(text) > 0 && text[0] == '@' {
		content, err := os.ReadFile(text[1:])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		text = string(content)
	}

	var result string
	var err error
	switch mode {
	case "encrypt":
		result, err = cipher.Encrypt(text)
	case "decrypt":
		result, err = cipher.Decrypt(text)
	default:
		fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Result:", result)
}
