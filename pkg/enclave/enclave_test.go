package enclave

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnclaveInitialization(t *testing.T) {
	// Initialize the secure enclave
	keyStore, err := InitializeEnclave()

	// Use Testify to assert the initialization
	assert.NoError(t, err, "Enclave initialization should not return an error")
	assert.NotNil(t, keyStore, "KeyStore should be initialized")

	// Check if the AES key was generated
	assert.NotNil(t, keyStore.AESKey, "AES key should be generated")
	assert.Len(t, keyStore.AESKey, keySize, "AES key should have the correct size")

	// Check RSA full and partial keys
	assert.NotNil(t, keyStore.RSAFullKey, "RSA full key should be generated")
	assert.Len(t, keyStore.RSAFullKey, rsaKeySize, "RSA full key should have the correct size")
	assert.NotNil(t, keyStore.RSAPartial, "RSA partial key should be generated")

	// Check ECDSA full and partial keys
	assert.NotNil(t, keyStore.ECDSAFull, "ECDSA full key should be generated")
	assert.Len(t, keyStore.ECDSAFull, rsaKeySize, "ECDSA full key should have the correct size")
	assert.NotNil(t, keyStore.ECDSAPartial, "ECDSA partial key should be generated")

	// Check Ed25519 full and partial keys
	assert.NotNil(t, keyStore.Ed25519Full, "Ed25519 full key should be generated")
	assert.Len(t, keyStore.Ed25519Full, rsaKeySize, "Ed25519 full key should have the correct size")
	assert.NotNil(t, keyStore.Ed25519Partial, "Ed25519 partial key should be generated")
}

func TestRSAOperations(t *testing.T) {
	keyStore, err := InitializeEnclave()
	assert.NoError(t, err, "Enclave initialization should not return an error")

	message := []byte("Test message for signing.")

	// Test RSA full signing
	signature, err := RSASign(message, keyStore)
	assert.NoError(t, err, "RSA signature operation should succeed")
	assert.NotNil(t, signature, "RSA signature should be generated")

	// Test RSA partial signing
	partialSig, err := RSAPartialSign(message, keyStore)
	assert.NoError(t, err, "RSA partial signature operation should succeed")
	assert.NotNil(t, partialSig, "RSA partial signature should be generated")
}

func TestECDSAOperations(t *testing.T) {
	keyStore, err := InitializeEnclave()
	assert.NoError(t, err, "Enclave initialization should not return an error")

	message := []byte("Test message for signing.")

	// Test ECDSA full signing
	signature, err := ECDSASign(message, keyStore)
	assert.NoError(t, err, "ECDSA signature operation should succeed")
	assert.NotNil(t, signature, "ECDSA signature should be generated")

	// Test ECDSA partial signing
	partialSig, err := ECDSAPartialSign(message, keyStore)
	assert.NoError(t, err, "ECDSA partial signature operation should succeed")
	assert.NotNil(t, partialSig, "ECDSA partial signature should be generated")
}

func TestEd25519Operations(t *testing.T) {
	keyStore, err := InitializeEnclave()
	assert.NoError(t, err, "Enclave initialization should not return an error")

	message := []byte("Test message for signing.")

	// Test Ed25519 full signing
	signature, err := Ed25519Sign(message, keyStore)
	assert.NoError(t, err, "Ed25519 signature operation should succeed")
	assert.NotNil(t, signature, "Ed25519 signature should be generated")

	// Test Ed25519 partial signing
	partialSig, err := Ed25519PartialSign(message, keyStore)
	assert.NoError(t, err, "Ed25519 partial signature operation should succeed")
	assert.NotNil(t, partialSig, "Ed25519 partial signature should be generated")
}

func TestAESOperations(t *testing.T) {
	keyStore, err := InitializeEnclave()
	assert.NoError(t, err, "Enclave initialization should not return an error")

	plaintext := []byte("Test data for AES encryption.")

	// Test AES encryption
	ciphertext, err := AESEncrypt(plaintext, keyStore)
	assert.NoError(t, err, "AES encryption operation should succeed")
	assert.NotNil(t, ciphertext, "AES ciphertext should be generated")

	// Test AES decryption
	decrypted, err := AESDecrypt(ciphertext, keyStore)
	assert.NoError(t, err, "AES decryption operation should succeed")
	assert.Equal(t, plaintext, decrypted, "Decrypted data should match the original plaintext")
}
