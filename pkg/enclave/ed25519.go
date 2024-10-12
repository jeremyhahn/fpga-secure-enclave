package enclave

import (
	"crypto/rand"
	"fmt"

	"github.com/hashicorp/vault/shamir"
	"github.com/jeremyhahn/fpga-secure-enclave/pkg/fpga"
)

// InitializeEd25519Key generates a random Ed25519 key, splits it using Shamir Secret Sharing, and loads both full and partial keys into FPGA
func InitializeEd25519Key(mappedMem []byte) ([]byte, []byte, error) {
	// Generate a random Ed25519 key (256 bytes for Ed25519)
	ed25519FullKey := make([]byte, rsaKeySize)
	_, err := rand.Read(ed25519FullKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate Ed25519 key: %v", err)
	}

	// Split the Ed25519 key using Shamir Secret Sharing
	shares, err := shamir.Split(ed25519FullKey, numShares, threshold)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to split Ed25519 key using Shamir: %v", err)
	}

	// Use the first share as the partial key shard
	ed25519PartialKey := shares[0]

	// Load the full Ed25519 key into the FPGA
	err = fpga.LoadKeyToFPGA(ed25519FullKey, 0x4000, mappedMem)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load Ed25519 full key to FPGA: %v", err)
	}

	// Load the partial Ed25519 key into the FPGA
	err = fpga.LoadKeyToFPGA(ed25519PartialKey, 0x4100, mappedMem)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load Ed25519 partial key to FPGA: %v", err)
	}

	fmt.Println("Ed25519 full and partial keys successfully loaded into the FPGA")
	return ed25519FullKey, ed25519PartialKey, nil
}

// Ed25519Sign performs a full Ed25519 signature using the complete private key
func Ed25519Sign(message []byte, keyStore *EnclaveKeyStore) ([]byte, error) {
	// Load full Ed25519 private key into FPGA
	err := fpga.LoadKeyToFPGA(keyStore.Ed25519Full, 0x4000, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load full Ed25519 key: %v", err)
	}

	// Perform Ed25519 signing (FPGA call to real signing logic would be added here)
	signature := make([]byte, 256) // Placeholder for actual signature
	copy(signature, message)       // Just copying message for testing

	fmt.Println("Performing Ed25519 full signing")
	return signature, nil
}

// Ed25519PartialSign performs a partial Ed25519 signature using a key shard
func Ed25519PartialSign(message []byte, keyStore *EnclaveKeyStore) ([]byte, error) {
	// Load Ed25519 partial key shard into FPGA
	err := fpga.LoadKeyToFPGA(keyStore.Ed25519Partial, 0x4100, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load Ed25519 partial key: %v", err)
	}

	// Perform Ed25519 partial signing (FPGA call to real signing logic would be added here)
	partialSignature := make([]byte, 256) // Placeholder for partial signature
	copy(partialSignature, message)       // Just copying message for testing

	fmt.Println("Performing Ed25519 partial signing")
	return partialSignature, nil
}
