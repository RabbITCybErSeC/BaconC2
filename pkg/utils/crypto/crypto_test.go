package crypto

import (
	"bytes"
	"testing"
)

func TestAESGCMEncryptDecrypt(t *testing.T) {
	key := RandomSymmetricKey()
	cipher, err := NewAESGCMCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher: %v", err)
	}

	plaintext := []byte("Hello, World! This is a test message.")
	
	ciphertext, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := cipher.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text doesn't match original.\nExpected: %s\nGot: %s", plaintext, decrypted)
	}
}

func TestAESCBCEncryptDecrypt(t *testing.T) {
	key := RandomSymmetricKey()
	cipher, err := NewAESCBCCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher: %v", err)
	}

	plaintext := []byte("Hello, World! This is a test message.")
	
	ciphertext, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := cipher.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text doesn't match original.\nExpected: %s\nGot: %s", plaintext, decrypted)
	}
}

func TestCipherContext(t *testing.T) {
	key := RandomSymmetricKey()
	cipher, err := NewAESGCMCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher: %v", err)
	}

	ctx := NewSymmetricCipherContext(cipher)
	plaintext := []byte("Test message for replay detection")

	// First encryption and decryption should work
	ciphertext, err := ctx.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := ctx.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text doesn't match original")
	}

	// Second decryption of the same ciphertext should fail (replay attack)
	_, err = ctx.Decrypt(ciphertext)
	if err != ErrReplayAttack {
		t.Errorf("Expected replay attack error, got: %v", err)
	}
}

func TestKeyFromBytes(t *testing.T) {
	// Test valid key
	validKey := make([]byte, AESKeySize)
	for i := range validKey {
		validKey[i] = byte(i)
	}

	key, err := KeyFromBytes(validKey)
	if err != nil {
		t.Fatalf("KeyFromBytes failed for valid key: %v", err)
	}

	if !bytes.Equal(key[:], validKey) {
		t.Errorf("Key doesn't match input")
	}

	// Test invalid key length
	invalidKey := make([]byte, 16)
	_, err = KeyFromBytes(invalidKey)
	if err != ErrInvalidKeyLength {
		t.Errorf("Expected ErrInvalidKeyLength, got: %v", err)
	}
}

func TestDeriveKeyFromPassword(t *testing.T) {
	password := "test-password-123"
	key1 := DeriveKeyFromPassword(password)
	key2 := DeriveKeyFromPassword(password)

	// Same password should produce same key
	if !bytes.Equal(key1[:], key2[:]) {
		t.Errorf("Same password produced different keys")
	}

	// Different password should produce different key
	key3 := DeriveKeyFromPassword("different-password")
	if bytes.Equal(key1[:], key3[:]) {
		t.Errorf("Different passwords produced same key")
	}
}


func BenchmarkAESGCMEncrypt(b *testing.B) {
	key := RandomSymmetricKey()
	cipher, _ := NewAESGCMCipher(key)
	plaintext := []byte("Benchmark test message for encryption performance")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cipher.Encrypt(plaintext)
	}
}

func BenchmarkAESGCMDecrypt(b *testing.B) {
	key := RandomSymmetricKey()
	cipher, _ := NewAESGCMCipher(key)
	plaintext := []byte("Benchmark test message for decryption performance")
	ciphertext, _ := cipher.Encrypt(plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cipher.Decrypt(ciphertext)
	}
}
