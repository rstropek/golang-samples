# Caesar Cipher Exercise in Go

## Objective

Your task is to implement the [Caesar cipher](https://cryptii.com/pipes/caesar-cipher) encryption and decryption in Go, using the provided interfaces. Additionally, you will create a command-line tool to demonstrate the functionality of your implementation.

## Interfaces

You are provided with two interfaces:

```go
type Encryptor interface {
  Encrypt(plain string) (string, error)
}

type Decryptor interface {
  Decrypt(cypher string) (string, error)
}
```

## Requirements

1. **Caesar Cipher Struct**:
   - Implement a struct (e.g., `CaesarCipher`) that will be used to perform the encryption and decryption.
   - The struct must have a constructor function that accepts an integer representing the shift value.
   - Ensure that your struct correctly implements both `Encryptor` and `Decryptor` interfaces.

2. **Encryption and Decryption Logic**:
   - The `Encrypt` method should apply the Caesar cipher technique to encrypt the provided plaintext.
   - The `Decrypt` method should reverse the encryption process to return the original text from the encrypted string.
   - Remember to handle both uppercase and lowercase letters and to leave non-alphabetic characters unchanged.
   - Ensure your code gracefully handles errors and edge cases, such as invalid shift values.

3. **Command-Line Tool**:
   - Create a command-line executable that uses your `CaesarCipher` struct.
   - The tool should accept three command-line arguments:
     - An integer representing the shift value.
     - A string indicating the mode: either `encrypt` or `decrypt`.
     - The text to be processed (either plaintext for encryption or ciphertext for decryption).
   - The tool should output the result of the encryption or decryption to the console.

## Bonus Challenge (Optional)

- Add [unit tests](https://gobyexample.com/testing) for your `CaesarCipher` struct to validate its functionality.
- Implement an additional feature in your command-line tool, such as reading input from a file or writing output to a file.
