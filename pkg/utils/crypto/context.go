package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"sync"
)

type SymmetricCipherContext struct {
	cipher ISymmetricCipher
	replay *sync.Map
}

func NewSymmetricCipherContext(cipher ISymmetricCipher) *SymmetricCipherContext {
	return &SymmetricCipherContext{
		cipher: cipher,
		replay: &sync.Map{},
	}
}

func (c *SymmetricCipherContext) Encrypt(plaintext []byte) ([]byte, error) {
	ciphertext, err := c.cipher.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) > 0 {
		digest := sha256.Sum256(ciphertext)
		b64Digest := base64.RawStdEncoding.EncodeToString(digest[:])
		c.replay.Store(b64Digest, true)
	}

	return ciphertext, nil
}

func (c *SymmetricCipherContext) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) > 0 {
		digest := sha256.Sum256(ciphertext)
		b64Digest := base64.RawStdEncoding.EncodeToString(digest[:])
		if _, ok := c.replay.LoadOrStore(b64Digest, true); ok {
			return nil, ErrReplayAttack
		}
	}

	plaintext, err := c.cipher.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (c *SymmetricCipherContext) Reset() {
	c.replay = &sync.Map{}
}
