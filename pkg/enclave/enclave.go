package enclave

import (
	"fmt"
	"os"
	"syscall"
)

const (
	keySize     = 32  // AES-256 key size in bytes
	rsaKeySize  = 256 // RSA key size in bytes
	numShares   = 5   // Number of shares for secret sharing
	threshold   = 3   // Threshold for secret sharing
	axiBaseAddr = 0xA0000000
	pageSize    = 4096
)

// EnclaveKeyStore holds the keys for AES, RSA, ECDSA, and Ed25519
type EnclaveKeyStore struct {
	AESKey         []byte
	RSAFullKey     []byte
	RSAPartial     []byte
	ECDSAFull      []byte
	ECDSAPartial   []byte
	Ed25519Full    []byte
	Ed25519Partial []byte
}

// Initialize the secure enclave by calling each key initializer
func InitializeEnclave() (*EnclaveKeyStore, error) {
	// Map memory for loading keys into FPGA
	memFile, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open /dev/mem: %v", err)
	}
	defer memFile.Close()

	mappedMem, err := syscall.Mmap(int(memFile.Fd()), axiBaseAddr, pageSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("failed to memory-map the AXI address region: %v", err)
	}
	defer syscall.Munmap(mappedMem)

	// Load AES key
	aesKey, err := InitializeAESKey(mappedMem)
	if err != nil {
		return nil, err
	}

	// Load RSA full and partial keys
	rsaFullKey, rsaPartial, err := InitializeRSAKey(mappedMem)
	if err != nil {
		return nil, err
	}

	// Load ECDSA full and partial keys
	ecdsaFullKey, ecdsaPartial, err := InitializeECDSAKey(mappedMem)
	if err != nil {
		return nil, err
	}

	// Load Ed25519 full and partial keys
	ed25519FullKey, ed25519Partial, err := InitializeEd25519Key(mappedMem)
	if err != nil {
		return nil, err
	}

	// Return the initialized EnclaveKeyStore
	return &EnclaveKeyStore{
		AESKey:         aesKey,
		RSAFullKey:     rsaFullKey,
		RSAPartial:     rsaPartial,
		ECDSAFull:      ecdsaFullKey,
		ECDSAPartial:   ecdsaPartial,
		Ed25519Full:    ed25519FullKey,
		Ed25519Partial: ed25519Partial,
	}, nil
}
