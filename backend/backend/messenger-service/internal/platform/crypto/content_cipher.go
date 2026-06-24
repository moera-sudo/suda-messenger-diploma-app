package crypto

import (
	"strings"

	"github.com/rs/zerolog/log"
)

// contentSentinelV1 marks a stored message body as AES-256-GCM ciphertext
// (key version 1). Strings without this prefix are treated as legacy plaintext,
// so old messages keep working without any backfill.
const contentSentinelV1 = "enc:v1:"

// ContentCipher encrypts/decrypts message bodies at the repository boundary.
// Implementations are safe to call on every read/write; the no-op variant lets
// the service run without an encryption key (dev) without breaking anything.
type ContentCipher interface {
	// EncryptContent returns the sentinel-prefixed ciphertext for plaintext.
	EncryptContent(plaintext string) string
	// DecryptContent returns plaintext: decrypts sentinel-prefixed values,
	// passes everything else (legacy plaintext) through unchanged.
	DecryptContent(stored string) string
}

type aesContentCipher struct {
	enc *Encryptor
}

func (c *aesContentCipher) EncryptContent(plaintext string) string {
	blob, err := c.enc.Encrypt([]byte(plaintext))
	if err != nil {
		// Never lose a message over an encryption hiccup — store plaintext.
		log.Error().Err(err).Msg("content encrypt failed, storing plaintext")
		return plaintext
	}
	return contentSentinelV1 + blob
}

func (c *aesContentCipher) DecryptContent(stored string) string {
	if !strings.HasPrefix(stored, contentSentinelV1) {
		return stored // legacy plaintext
	}
	plain, err := c.enc.Decrypt(1, strings.TrimPrefix(stored, contentSentinelV1))
	if err != nil {
		log.Error().Err(err).Msg("content decrypt failed")
		return stored
	}
	return string(plain)
}

// noopContentCipher is used when no key is configured: identity transforms.
type noopContentCipher struct{}

func (noopContentCipher) EncryptContent(plaintext string) string { return plaintext }
func (noopContentCipher) DecryptContent(stored string) string    { return stored }

// NewContentCipher builds a ContentCipher from a hex key. An empty key yields a
// no-op cipher (encryption disabled) so the service still runs.
func NewContentCipher(keyHex string) (ContentCipher, error) {
	if strings.TrimSpace(keyHex) == "" {
		return noopContentCipher{}, nil
	}
	enc, err := NewEncryptor(1, keyHex)
	if err != nil {
		return nil, err
	}
	return &aesContentCipher{enc: enc}, nil
}

// Enabled reports whether the cipher actually encrypts (true) or is a no-op.
func Enabled(c ContentCipher) bool {
	_, noop := c.(noopContentCipher)
	return !noop
}
