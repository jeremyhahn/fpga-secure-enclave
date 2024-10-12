package enclave

import (
	"crypto/rand"
	"fmt"

	"github.com/hashicorp/vault/shamir"
	"github.com/jeremyhahn/fpga-secure-enclave/pkg/fpga"
)

// InitializeRSAKey generates a random RSA key, splits it using Shamir Secret Sharing, and loads both full and partial keys into FPGA
func InitializeRSAKey(mappedMem []byte) ([]byte, []byte, error) {
	// Generate a random RSA key (256 bytes for 2048-bit RSA)
	rsaFullKey := make([]byte, rsaKeySize)
	_, err := rand.Read(rsaFullKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key: %v", err)
	}

	// Split the RSA key using Shamir Secret Sharing
	// n = total shares, threshold = minimum number of shares to reassemble the key
	shares, err := shamir.Split(rsaFullKey, numShares, threshold)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to split RSA key using Shamir: %v", err)
	}

	// Use the first share as the partial key shard
	rsaPartialKey := shares[0]

	// Load the full RSA key into the FPGA
	err = fpga.LoadKeyToFPGA(rsaFullKey, 0x2000, mappedMem)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load RSA full key to FPGA: %v", err)
	}

	// Load the partial RSA key into the FPGA
	err = fpga.LoadKeyToFPGA(rsaPartialKey, 0x2100, mappedMem)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load RSA partial key to FPGA: %v", err)
	}

	fmt.Println("RSA full and partial keys successfully loaded into the FPGA")
	return rsaFullKey, rsaPartialKey, nil
}

// RSASign performs a full RSA signature using the complete private key
func RSASign(message []byte, keyStore *EnclaveKeyStore) ([]byte, error) {
	// Load full RSA private key into FPGA
	err := fpga.LoadKeyToFPGA(keyStore.RSAFullKey, 0x2000, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load full RSA key: %v", err)
	}

	// Simulate RSA signing operation (FPGA call would be added here)
	signature := make([]byte, 256) // Placeholder for actual signature
	copy(signature, message)       // Just copying message for testing

	fmt.Println("Performing RSA full signing")
	return signature, nil
}

// RSAPartialSign performs a partial RSA signature using a key shard (threshold signing)
func RSAPartialSign(message []byte, keyStore *EnclaveKeyStore) ([]byte, error) {
	// Load RSA partial key shard into FPGA
	err := fpga.LoadKeyToFPGA(keyStore.RSAPartial, 0x2100, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load RSA partial key: %v", err)
	}

	// Simulate RSA partial signing operation (FPGA call would be added here)
	partialSignature := make([]byte, 256) // Placeholder for partial signature
	copy(partialSignature, message)       // Just copying message for testing

	fmt.Println("Performing RSA partial signing")
	return partialSignature, nil
}
