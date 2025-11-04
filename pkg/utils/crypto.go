package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rc4"
	"encoding/binary"
	"errors"
	"io"
	"math"
)

func RC4Encrypt(data []byte, key []byte) []byte {
	c, err := rc4.NewCipher(key)
	if err != nil {
		return make([]byte, 0)
	}
	cipherText := make([]byte, len(data))
	c.XORKeyStream(cipherText, data)
	return cipherText
}

func PreludeEncrypt(data []byte, key []byte, iv []byte) []byte {
	plainText, err := pad(data, aes.BlockSize)
	if err != nil {
		return make([]byte, 0)
	}
	block, _ := aes.NewCipher(key)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	if len(iv) == 0 {
		iv = cipherText[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return make([]byte, 0)
		}
	} else {
		copy(cipherText[:aes.BlockSize], iv)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)
	return cipherText
}

func PreludeDecrypt(data []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)
	data, _ = unpad(data, aes.BlockSize)
	return data
}

func pad(buf []byte, size int) ([]byte, error) {
	bufLen := len(buf)
	padLen := size - bufLen%size
	padded := make([]byte, bufLen+padLen)
	copy(padded, buf)
	for i := 0; i < padLen; i++ {
		padded[bufLen+i] = byte(padLen)
	}
	return padded, nil
}

func unpad(padded []byte, size int) ([]byte, error) {
	if len(padded)%size != 0 {
		return nil, errors.New("pkcs7: Padded value wasn't in correct size")
	}
	bufLen := len(padded) - int(padded[len(padded)-1])
	buf := make([]byte, bufLen)
	copy(buf, padded[:bufLen])
	return buf, nil
}

func Intn(n int) int {
	if n <= 0 {
		panic("secure.Intn: non-positive n")
	}
	un := uint64(n)

	limit := (math.MaxUint64 / un) * un

	for {
		x := mustRandUint64()
		if x < limit {
			return int(x % un)
		}
	}
}

func Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("secure.Shuffle: negative n")
	}
	for i := n - 1; i > 0; i-- {
		j := Intn(i + 1)
		if i != j {
			swap(i, j)
		}
	}
}

func mustRandUint64() uint64 {
	var b [8]byte
	if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
		panic("secure: crypto/rand failure: " + err.Error())
	}
	return binary.LittleEndian.Uint64(b[:])
}

// Int63n returns a uniform int64 in [0, n).
// Panics if n <= 0 or if crypto/rand fails.
func Int63n(n int64) int64 {
	if n <= 0 {
		panic("secure.Int63n: non-positive n")
	}
	un := uint64(n)

	const max63 = uint64(1<<63 - 1)
	limit := max63 - (max63+1)%un

	for {
		x := randUint63()
		if x <= limit {
			return int64(x % un)
		}
	}
}

func randUint63() uint64 {
	var b [8]byte
	if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
		panic("secure: crypto/rand failure: " + err.Error())
	}
	return binary.LittleEndian.Uint64(b[:]) & (1<<63 - 1)
}

func Float64() float64 {
	u := randUint53()
	const inv1p53 = 1.0 / (1 << 53)
	return float64(u) * inv1p53
}

func randUint53() uint64 {
	var b [8]byte
	if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
		panic("secure: crypto/rand failure: " + err.Error())
	}
	return binary.LittleEndian.Uint64(b[:]) >> 11
}
