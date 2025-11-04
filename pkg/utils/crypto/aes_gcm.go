package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

const (
	AESKeySize      = 32
	AESGCMNonceSize = 12
)

type AESGCMCipher struct {
	key  [AESKeySize]byte
	aead cipher.AEAD
}

func NewAESGCMCipher(key [AESKeySize]byte) (*AESGCMCipher, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &AESGCMCipher{
		key:  key,
		aead: aead,
	}, nil
}

func NewAESGCMCipherFromBytes(keyBytes []byte) (*AESGCMCipher, error) {
	key, err := KeyFromBytes(keyBytes)
	if err != nil {
		return nil, err
	}
	return NewAESGCMCipher(key)
}

func (c *AESGCMCipher) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := c.aead.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (c *AESGCMCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := c.aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrInvalidCiphertext
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	return plaintext, nil
}

func (c *AESGCMCipher) Key() []byte {
	return c.key[:]
}

type AESCBCCipher struct {
	key [AESKeySize]byte
}

func NewAESCBCCipher(key [AESKeySize]byte) (*AESCBCCipher, error) {
	return &AESCBCCipher{
		key: key,
	}, nil
}

func NewAESCBCCipherFromBytes(keyBytes []byte) (*AESCBCCipher, error) {
	key, err := KeyFromBytes(keyBytes)
	if err != nil {
		return nil, err
	}
	return NewAESCBCCipher(key)
}

func (c *AESCBCCipher) Encrypt(plaintext []byte) ([]byte, error) {
	padded, err := pkcs7Pad(plaintext, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(c.key[:])
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(padded))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], padded)

	return ciphertext, nil
}

func (c *AESCBCCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(c.key[:])
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, ErrInvalidCiphertext
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext, err = pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	return plaintext, nil
}

func (c *AESCBCCipher) Key() []byte {
	return c.key[:]
}

func pkcs7Pad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 || blockSize > 255 {
		return nil, ErrInvalidKeyLength
	}
	padLen := blockSize - len(data)%blockSize
	padding := make([]byte, padLen)
	for i := range padding {
		padding[i] = byte(padLen)
	}
	return append(data, padding...), nil
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, ErrInvalidCiphertext
	}
	if len(data)%blockSize != 0 {
		return nil, ErrInvalidCiphertext
	}
	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > blockSize {
		return nil, ErrInvalidCiphertext
	}
	if padLen > len(data) {
		return nil, ErrInvalidCiphertext
	}
	for i := 0; i < padLen; i++ {
		if data[len(data)-1-i] != byte(padLen) {
			return nil, ErrInvalidCiphertext
		}
	}
	return data[:len(data)-padLen], nil
}

func deriveKeyFrom(data []byte) [AESKeySize]byte {
	digest := sha256.Sum256(data)
	var key [AESKeySize]byte
	copy(key[:], digest[:AESKeySize])
	return key
}

func KeyFromBytes(data []byte) ([AESKeySize]byte, error) {
	var key [AESKeySize]byte
	if len(data) != AESKeySize {
		return RandomSymmetricKey(), ErrInvalidKeyLength
	}
	copy(key[:], data)
	return key, nil
}

func RandomSymmetricKey() [AESKeySize]byte {
	randBuf := make([]byte, 64)
	if _, err := rand.Read(randBuf); err != nil {
		panic(err)
	}
	return deriveKeyFrom(randBuf)
}

func DeriveKeyFromPassword(password string) [AESKeySize]byte {
	return deriveKeyFrom([]byte(password))
}
