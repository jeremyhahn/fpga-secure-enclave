package enclave

import (
	"fmt"
	"log"

	"github.com/jeremyhahn/fpga-secure-enclave/pkg/enclave"
)

func main() {
	// Initialize the secure enclave
	keyStore, err := enclave.InitializeEnclave()
	if err != nil {
		log.Fatalf("Failed to initialize enclave: %v", err)
	}

	message := []byte("Test message for signing.")

	// Perform RSA signing operations
	rsaSignature, err := enclave.RSASign(message, keyStore)
	if err != nil {
		log.Fatalf("RSA full signing failed: %v", err)
	}
	fmt.Printf("RSA Full Signature: %x\n", rsaSignature)

	rsaPartialSignature, err := enclave.RSAPartialSign(message, keyStore)
	if err != nil {
		log.Fatalf("RSA partial signing failed: %v", err)
	}
	fmt.Printf("RSA Partial Signature: %x\n", rsaPartialSignature)

	// Perform ECDSA signing operations
	ecdsaSignature, err := enclave.ECDSASign(message, keyStore)
	if err != nil {
		log.Fatalf("ECDSA full signing failed: %v", err)
	}
	fmt.Printf("ECDSA Full Signature: %x\n", ecdsaSignature)

	ecdsaPartialSignature, err := enclave.ECDSAPartialSign(message, keyStore)
	if err != nil {
		log.Fatalf("ECDSA partial signing failed: %v", err)
	}
	fmt.Printf("ECDSA Partial Signature: %x\n", ecdsaPartialSignature)

	// Perform Ed25519 signing operations
	ed25519Signature, err := enclave.Ed25519Sign(message, keyStore)
	if err != nil {
		log.Fatalf("Ed25519 full signing failed: %v", err)
	}
	fmt.Printf("Ed25519 Full Signature: %x\n", ed25519Signature)

	ed25519PartialSignature, err := enclave.Ed25519PartialSign(message, keyStore)
	if err != nil {
		log.Fatalf("Ed25519 partial signing failed: %v", err)
	}
	fmt.Printf("Ed25519 Partial Signature: %x\n", ed25519PartialSignature)

	// Test AES encryption and decryption
	plaintext := []byte("Test data for AES encryption.")

	ciphertext, err := enclave.AESEncrypt(plaintext, keyStore)
	if err != nil {
		log.Fatalf("AES encryption failed: %v", err)
	}
	fmt.Printf("AES Ciphertext: %x\n", ciphertext)

	decryptedText, err := enclave.AESDecrypt(ciphertext, keyStore)
	if err != nil {
		log.Fatalf("AES decryption failed: %v", err)
	}
	fmt.Printf("Decrypted Text: %s\n", string(decryptedText))
}
