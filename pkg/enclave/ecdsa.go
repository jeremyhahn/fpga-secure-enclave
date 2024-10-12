package enclave

import (
	"crypto/rand"
	"fmt"

	"github.com/hashicorp/vault/shamir"
	"github.com/jeremyhahn/fpga-secure-enclave/pkg/fpga"
)

// InitializeECDSAKey generates a random ECDSA key, splits it using Shamir Secret Sharing, and loads both full and partial keys into FPGA
func InitializeECDSAKey(mappedMem []byte) ([]byte, []byte, error) {
	// Generate a random ECDSA key (256 bytes for 256-bit curve)
	ecdsaFullKey := make([]byte, rsaKeySize)
	_, err := rand.Read(ecdsaFullKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate ECDSA key: %v", err)
	}

	// Split the ECDSA key using Shamir Secret Sharing
	shares, err := shamir.Split(ecdsaFullKey, numShares, threshold)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to split ECDSA key using Shamir: %v", err)
	}

	// Use the first share as the partial key shard
	ecdsaPartialKey := shares[0]

	// Load the full ECDSA key into the FPGA
	err = fpga.LoadKeyToFPGA(ecdsaFullKey, 0x3000, mappedMem)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load ECDSA full key to FPGA: %v", err)
	}

	// Load the partial ECDSA key into the FPGA
	err = fpga.LoadKeyToFPGA(ecdsaPartialKey, 0x3100, mappedMem)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load ECDSA partial key to FPGA: %v", err)
	}

	fmt.Println("ECDSA full and partial keys successfully loaded into the FPGA")
	return ecdsaFullKey, ecdsaPartialKey, nil
}

// ECDSASign performs a full ECDSA signature using the complete private key
func ECDSASign(message []byte, keyStore *EnclaveKeyStore) ([]byte, error) {
	// Load full ECDSA private key into FPGA
	err := fpga.LoadKeyToFPGA(keyStore.ECDSAFull, 0x3000, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load full ECDSA key: %v", err)
	}

	// Perform ECDSA signing (FPGA call to real signing logic would be added here)
	signature := make([]byte, 256) // Placeholder for actual signature
	copy(signature, message)       // Just copying message for testing

	fmt.Println("Performing ECDSA full signing")
	return signature, nil
}

// ECDSAPartialSign performs a partial ECDSA signature using a key shard
func ECDSAPartialSign(message []byte, keyStore *EnclaveKeyStore) ([]byte, error) {
	// Load ECDSA partial key shard into FPGA
	err := fpga.LoadKeyToFPGA(keyStore.ECDSAPartial, 0x3100, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load ECDSA partial key: %v", err)
	}

	// Perform ECDSA partial signing (FPGA call to real signing logic would be added here)
	partialSignature := make([]byte, 256) // Placeholder for partial signature
	copy(partialSignature, message)       // Just copying message for testing

	fmt.Println("Performing ECDSA partial signing")
	return partialSignature, nil
}
