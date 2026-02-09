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

type EncryptorDecryptor interface {
	Encryptor
	Decryptor
}

type caesarCipher struct {
	shift int
}

func NewCaesarCipher(shift int) EncryptorDecryptor {
	return &caesarCipher{shift: shift}
}

func (cc *caesarCipher) Encrypt(plain string) (string, error) {
	return cc.shiftText(plain, cc.shift)
}

func (cc *caesarCipher) Decrypt(cipher string) (string, error) {
	return cc.shiftText(cipher, -cc.shift)
}

func (cc *caesarCipher) shiftText(text string, shift int) (string, error) {
	shift = ((shift % 26) + 26) % 26

	// A small allow-list to demonstrate using a map.
	allowedChars := map[rune]bool{
		' ':  true,
		'.':  true,
		',':  true,
		'!':  true,
		'?':  true,
		'\'': true,
		'-':  true,
		'\n': true,
	}
	for d := '0'; d <= '9'; d++ {
		allowedChars[d] = true
	}

	var resultBuilder strings.Builder
	resultBuilder.Grow(len(text))
	for _, char := range text {
		if char >= 'a' && char <= 'z' {
			resultBuilder.WriteRune(((char-'a'+rune(shift))%26 + 'a'))
		} else if char >= 'A' && char <= 'Z' {
			resultBuilder.WriteRune(((char-'A'+rune(shift))%26 + 'A'))
		} else {
			if !allowedChars[char] {
				return "", fmt.Errorf("invalid character: %q", char)
			}

			resultBuilder.WriteRune(char)
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

	if text == "" {
		fmt.Fprintln(os.Stderr, "Error: missing text. Provide -text.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cipher := NewCaesarCipher(shift)

	// If text starts with @, read from file
	if len(text) > 0 && text[0] == '@' {
		content, err := os.ReadFile(text[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
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
		fmt.Fprintln(os.Stderr, "Invalid mode. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
