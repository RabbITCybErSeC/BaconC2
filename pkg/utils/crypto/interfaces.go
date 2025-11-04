package crypto

import "errors"

var (
	ErrInvalidKeyLength  = errors.New("invalid key length")
	ErrDecryptFailed     = errors.New("decryption failed")
	ErrEncryptFailed     = errors.New("encryption failed")
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	ErrReplayAttack      = errors.New("replay attack detected")
)

type ISymmetricCipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	Key() []byte
}

type IAsymmetricCipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	PublicKey() string
	PrivateKey() string
}

type ICipherContext interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	Reset()
}
