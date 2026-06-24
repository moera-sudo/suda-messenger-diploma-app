package crypto

import (
	"strings"
	"testing"
)

// 32-byte (64 hex char) test key.
const testKeyHex = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"

func TestContentCipher_RoundTrip(t *testing.T) {
	c, err := NewContentCipher(testKeyHex)
	if err != nil {
		t.Fatalf("NewContentCipher: %v", err)
	}
	if !Enabled(c) {
		t.Fatal("expected cipher to be enabled with a key")
	}

	plaintext := "Привет, мир! 🔐 hello"
	stored := c.EncryptContent(plaintext)

	if !strings.HasPrefix(stored, contentSentinelV1) {
		t.Fatalf("ciphertext missing sentinel prefix: %q", stored)
	}
	if strings.Contains(stored, plaintext) {
		t.Fatal("stored value should not contain plaintext")
	}
	if got := c.DecryptContent(stored); got != plaintext {
		t.Fatalf("round-trip mismatch: got %q want %q", got, plaintext)
	}
}

func TestContentCipher_LegacyPlaintextPassthrough(t *testing.T) {
	c, _ := NewContentCipher(testKeyHex)

	legacy := "old plaintext message without sentinel"
	if got := c.DecryptContent(legacy); got != legacy {
		t.Fatalf("legacy passthrough failed: got %q want %q", got, legacy)
	}
}

func TestContentCipher_DecryptFailsClosed(t *testing.T) {
	c, _ := NewContentCipher(testKeyHex)

	// Sentinel present but body is garbage -> return stored unchanged, no panic.
	bad := contentSentinelV1 + "not-valid-base64-or-cipher!!"
	if got := c.DecryptContent(bad); got != bad {
		t.Fatalf("expected unchanged on decrypt failure, got %q", got)
	}
}

func TestContentCipher_NoopWhenNoKey(t *testing.T) {
	c, err := NewContentCipher("")
	if err != nil {
		t.Fatalf("NewContentCipher(empty): %v", err)
	}
	if Enabled(c) {
		t.Fatal("expected disabled (no-op) cipher with empty key")
	}

	p := "plaintext"
	if c.EncryptContent(p) != p {
		t.Fatal("noop EncryptContent should be identity")
	}
	if c.DecryptContent(p) != p {
		t.Fatal("noop DecryptContent should be identity")
	}
}

func TestContentCipher_NonceIsRandom(t *testing.T) {
	c, _ := NewContentCipher(testKeyHex)
	a := c.EncryptContent("same")
	b := c.EncryptContent("same")
	if a == b {
		t.Fatal("two encryptions of the same plaintext must differ (random nonce)")
	}
	if c.DecryptContent(a) != "same" || c.DecryptContent(b) != "same" {
		t.Fatal("both ciphertexts must decrypt to the same plaintext")
	}
}

func TestNewContentCipher_BadKey(t *testing.T) {
	if _, err := NewContentCipher("zzzz"); err == nil {
		t.Fatal("expected error for invalid hex key")
	}
	if _, err := NewContentCipher("00112233"); err == nil {
		t.Fatal("expected error for short key")
	}
}
