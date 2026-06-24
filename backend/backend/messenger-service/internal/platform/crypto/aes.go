// Package crypto provides AES-256-GCM symmetric encryption for sensitive data.
// In messenger-service it is used to encrypt message content at rest
// (messenger_messages.content) for DIRECT/GROUP chats.
//
// Format of an encrypted blob (after base64 decode):
//
//	[nonce (12 bytes) | ciphertext+tag (variable)]
//
// The auth tag is part of GCM's output (last 16 bytes). The key version is not
// embedded in the blob — for message content it is carried by a sentinel prefix
// on the stored string (see content_cipher.go). The Encryptor below mirrors the
// proven pattern from transaction-service.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Encryptor — AES-256-GCM encryptor with key-version support. The current
// version is used for new encryptions; all registered versions can decrypt
// (to support gradual key rotation).
type Encryptor struct {
	currentVersion uint8
	keys           map[uint8]cipher.AEAD
}

// NewEncryptor builds an encryptor from a single key version.
// keyHex — 64-char hex string (32 bytes = AES-256).
func NewEncryptor(currentVersion uint8, keyHex string) (*Encryptor, error) {
	aead, err := buildAEAD(keyHex)
	if err != nil {
		return nil, err
	}
	return &Encryptor{
		currentVersion: currentVersion,
		keys:           map[uint8]cipher.AEAD{currentVersion: aead},
	}, nil
}

// AddLegacyKey registers an older key version (decrypt-only). Not thread-safe;
// call before the service starts serving.
func (e *Encryptor) AddLegacyKey(version uint8, keyHex string) error {
	if _, exists := e.keys[version]; exists {
		return fmt.Errorf("crypto: key version %d already registered", version)
	}
	aead, err := buildAEAD(keyHex)
	if err != nil {
		return err
	}
	e.keys[version] = aead
	return nil
}

// CurrentVersion — version under which new data is encrypted.
func (e *Encryptor) CurrentVersion() uint8 { return e.currentVersion }

// Encrypt encrypts plaintext under the current key version. Returns a
// base64-encoded string (nonce || ciphertext+tag).
func (e *Encryptor) Encrypt(plaintext []byte) (string, error) {
	aead, ok := e.keys[e.currentVersion]
	if !ok {
		return "", fmt.Errorf("crypto: current key version %d not configured", e.currentVersion)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("crypto: generate nonce: %w", err)
	}

	ciphertext := aead.Seal(nil, nonce, plaintext, nil)

	blob := make([]byte, 0, len(nonce)+len(ciphertext))
	blob = append(blob, nonce...)
	blob = append(blob, ciphertext...)

	return base64.StdEncoding.EncodeToString(blob), nil
}

// Decrypt decrypts a blob that was encrypted with the given key version.
func (e *Encryptor) Decrypt(version uint8, encoded string) ([]byte, error) {
	aead, ok := e.keys[version]
	if !ok {
		return nil, fmt.Errorf("crypto: key version %d not configured", version)
	}

	blob, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("crypto: decode base64: %w", err)
	}

	nonceSize := aead.NonceSize()
	if len(blob) < nonceSize {
		return nil, errors.New("crypto: ciphertext too short")
	}

	nonce, ciphertext := blob[:nonceSize], blob[nonceSize:]

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// Don't leak the cause — could be tamper or wrong version.
		return nil, errors.New("crypto: decrypt failed")
	}

	return plaintext, nil
}

// buildAEAD builds AES-256-GCM from a hex key.
func buildAEAD(keyHex string) (cipher.AEAD, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("crypto: key hex decode: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("crypto: key must be 32 bytes (got %d)", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("crypto: new cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("crypto: new gcm: %w", err)
	}

	return aead, nil
}
