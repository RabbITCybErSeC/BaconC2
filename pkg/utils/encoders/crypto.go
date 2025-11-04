package encoders

import "github.com/RabbITCybErSeC/BaconC2/pkg/utils/crypto"

type CryptoEncoder struct {
	cipher crypto.ISymmetricCipher
}

func NewCryptoEncoder(cipher crypto.ISymmetricCipher) *CryptoEncoder {
	return &CryptoEncoder{
		cipher: cipher,
	}
}

func NewAESGCMEncoder(key [crypto.AESKeySize]byte) (*CryptoEncoder, error) {
	cipher, err := crypto.NewAESGCMCipher(key)
	if err != nil {
		return nil, err
	}
	return NewCryptoEncoder(cipher), nil
}

func NewAESCBCEncoder(key [crypto.AESKeySize]byte) (*CryptoEncoder, error) {
	cipher, err := crypto.NewAESCBCCipher(key)
	if err != nil {
		return nil, err
	}
	return NewCryptoEncoder(cipher), nil
}

func (e *CryptoEncoder) Encode(data []byte) ([]byte, error) {
	return e.cipher.Encrypt(data)
}

func (e *CryptoEncoder) Decode(data []byte) ([]byte, error) {
	return e.cipher.Decrypt(data)
}

func (e *CryptoEncoder) Cipher() crypto.ISymmetricCipher {
	return e.cipher
}
