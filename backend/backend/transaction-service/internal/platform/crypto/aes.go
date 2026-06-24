// Package crypto provides AES-256-GCM symmetric encryption for sensitive data,
// primarily users' EVM private keys stored in tx_wallets.encrypted_private_key.
//
// Format of an encrypted ciphertext (after base64 decode):
//
//	[nonce (12 bytes) | ciphertext+tag (variable)]
//
// The auth tag is part of GCM's output (last 16 bytes of ciphertext+tag).
// We don't store the key version in the ciphertext itself — it lives in
// tx_wallets.key_version. The Encryptor below knows about all historical
// key versions; the current version is used for new encryptions, older
// versions are only used for decryption of legacy data.
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

// Encryptor — AES-256-GCM шифровальщик с поддержкой версий ключа.
// Текущая версия используется для новых шифрований; все известные версии
// поддерживаются для дешифрования (для постепенной миграции при ротации).
type Encryptor struct {
	currentVersion uint8
	keys           map[uint8]cipher.AEAD
}

// NewEncryptor создаёт энкриптор с одной версией ключа.
// keyHex — 64-символьная hex-строка (32 байта = AES-256).
func NewEncryptor(currentVersion uint8, keyHex string) (*Encryptor, error) {
	aead, err := buildAEAD(keyHex)
	if err != nil {
		return nil, err
	}
	return &Encryptor{
		currentVersion: currentVersion,
		keys: map[uint8]cipher.AEAD{
			currentVersion: aead,
		},
	}, nil
}

// AddLegacyKey добавляет ключ старой версии. Используется только для
// дешифрования (Encrypt всегда пишет под currentVersion).
// Вызывается до начала работы сервиса, не thread-safe.
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

// CurrentVersion — версия ключа, под которой шифруются новые данные.
func (e *Encryptor) CurrentVersion() uint8 {
	return e.currentVersion
}

// Encrypt шифрует plaintext под текущей версией ключа. Возвращает
// base64-кодированную строку (nonce || ciphertext+tag).
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

	// Склеиваем nonce и ciphertext в один блоб, кодируем в base64.
	blob := make([]byte, 0, len(nonce)+len(ciphertext))
	blob = append(blob, nonce...)
	blob = append(blob, ciphertext...)

	return base64.StdEncoding.EncodeToString(blob), nil
}

// Decrypt дешифрует blob, зашифрованный указанной версией ключа.
// Версия читается из БД (tx_wallets.key_version) и передаётся вызывающим.
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
		// Не подсвечиваем причину — может быть и тампер, и неверная версия.
		return nil, errors.New("crypto: decrypt failed")
	}

	return plaintext, nil
}

// buildAEAD создаёт GCM-режим поверх AES-256 из hex-ключа.
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