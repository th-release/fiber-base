package utils

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

func EncodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func DecodeBase64(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// XorEncrypt XORs each byte of text with the key and returns the result as a hex string.
func XorEncrypt(text, key string) string {
	textBytes := []byte(text)
	keyBytes := []byte(key)
	var b strings.Builder
	for i, c := range textBytes {
		b.WriteString(hex.EncodeToString([]byte{c ^ keyBytes[i%len(keyBytes)]}))
	}
	return b.String()
}

// XorDecrypt decodes a hex string produced by XorEncrypt back to the original text.
func XorDecrypt(hexText, key string) (string, error) {
	keyBytes := []byte(key)
	raw, err := hex.DecodeString(hexText)
	if err != nil {
		return "", err
	}
	result := make([]byte, len(raw))
	for i, b := range raw {
		result[i] = b ^ keyBytes[i%len(keyBytes)]
	}
	return string(result), nil
}

func Sha512Hex(input string) string {
	hash := sha512.Sum512([]byte(input))
	return fmt.Sprintf("%x", hash)
}
