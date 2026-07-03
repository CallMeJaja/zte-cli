package router

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// SHA256Hex returns the hex-encoded SHA256 hash of the input string.
func SHA256Hex(text string) string {
	hash := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", hash)
}

// deriveKey derives a 32-byte AES key from the session token using SHA256.
func deriveKey(sessionToken string) []byte {
	hash := sha256.Sum256([]byte(sessionToken))
	return hash[:]
}

// deriveIV derives a 16-byte IV from a substring of the session token using SHA256.
// The browser uses session_token[5:15] as the IV source.
// CryptoJS internally truncates the SHA256 result to 16 bytes (first 4 words).
func deriveIV(sessionToken string) []byte {
	ivSource := sessionToken
	if len(sessionToken) > 15 {
		ivSource = sessionToken[5:15]
	}
	hash := sha256.Sum256([]byte(ivSource))
	return hash[:16] // Truncate to 16 bytes for AES-CBC IV
}

// Decrypt decrypts an AES-CBC ZeroPadding encrypted string (as used by CryptoJS in the router's web UI).
// The encrypted string is base64-encoded.
func Decrypt(encryptedBase64 string, sessionToken string) (string, error) {
	if encryptedBase64 == "" || sessionToken == "" {
		return encryptedBase64, nil
	}

	key := deriveKey(sessionToken)
	iv := deriveIV(sessionToken)

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes cipher creation failed: %w", err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove ZeroPadding (null bytes at the end)
	plaintext = removeZeroPadding(plaintext)

	return string(plaintext), nil
}

// Encrypt encrypts a plaintext string using AES-CBC ZeroPadding (CryptoJS compatible).
// Returns a base64-encoded ciphertext string.
func Encrypt(plaintext string, sessionToken string) (string, error) {
	if plaintext == "" || sessionToken == "" {
		return plaintext, nil
	}

	key := deriveKey(sessionToken)
	iv := deriveIV(sessionToken)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes cipher creation failed: %w", err)
	}

	// Apply ZeroPadding
	padded := addZeroPadding([]byte(plaintext), aes.BlockSize)

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(padded))
	mode.CryptBlocks(ciphertext, padded)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// removeZeroPadding removes trailing null bytes (ZeroPadding used by CryptoJS).
func removeZeroPadding(data []byte) []byte {
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] != 0 {
			return data[:i+1]
		}
	}
	return data
}

// addZeroPadding adds null bytes to align plaintext to the AES block size.
func addZeroPadding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	if padding == blockSize {
		return data
	}
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	return padded
}
