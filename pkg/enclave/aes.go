package enclave

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/jeremyhahn/fpga-secure-enclave/pkg/fpga"
)

// InitializeAESKey generates a random AES-256 key and loads it into the FPGA
func InitializeAESKey(mappedMem []byte) ([]byte, error) {
	// Generate a random AES-256 key (32 bytes)
	aesKey := make([]byte, keySize)
	_, err := rand.Read(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %v", err)
	}

	// Load the AES key into FPGA memory using the AXI interface
	err = fpga.LoadKeyToFPGA(aesKey, 0x1000, mappedMem)
	if err != nil {
		return nil, fmt.Errorf("failed to load AES key to FPGA: %v", err)
	}

	fmt.Println("AES-256 key successfully loaded into the FPGA")
	return aesKey, nil
}

// AESEncrypt encrypts data using AES-256 in CTR mode
func AESEncrypt(plaintext []byte, keyStore *EnclaveKeyStore) ([]byte, error) {
	// Load AES key into FPGA
	err := fpga.LoadKeyToFPGA(keyStore.AESKey, 0x1000, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load AES key: %v", err)
	}

	// Perform AES encryption (FPGA call to real encryption logic would be added here)
	ciphertext := make([]byte, len(plaintext)) // Placeholder for encrypted data
	copy(ciphertext, plaintext)                // Just copying for now

	fmt.Println("Performing AES encryption")
	return ciphertext, nil
}

// AESDecrypt decrypts data using AES-256 in CTR mode
func AESDecrypt(ciphertext []byte, keyStore *EnclaveKeyStore) ([]byte, error) {
	// Load AES key into FPGA
	err := fpga.LoadKeyToFPGA(keyStore.AESKey, 0x1000, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load AES key: %v", err)
	}

	// Perform AES decryption (FPGA call to real decryption logic would be added here)
	plaintext := make([]byte, len(ciphertext)) // Placeholder for decrypted data
	copy(plaintext, ciphertext)                // Just copying for now

	fmt.Println("Performing AES decryption")
	return plaintext, nil
}

// EncryptCodeAES encrypts the code using AES-256 in CTR mode
func EncryptCodeAES(code []byte, key []byte) ([]byte, []byte, error) {
	// Create AES block cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	// Encrypt the code using AES in CTR mode
	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(code))
	stream.XORKeyStream(ciphertext, code)

	return ciphertext, iv, nil
}
