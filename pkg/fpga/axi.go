package fpga

import (
	"fmt"
	"syscall"
	"unsafe"
)

// LoadKeyToFPGA loads a key into the FPGA memory via AXI
func LoadKeyToFPGA(key []byte, axiOffset uint32, mappedMem []byte) error {

	// Get the page size once and store it in a variable
	pageSize := syscall.Getpagesize()

	// Ensure the key fits in the mapped memory
	if len(key)*4 > pageSize {
		return fmt.Errorf("key size exceeds mapped memory size")
	}

	// Access the memory mapped region and load the key into the FPGA memory
	for i := 0; i < len(key); i++ {
		axiMem := (*[1 << 30]uint32)(unsafe.Pointer(&mappedMem[0])) // Using a large array length for maximum compatibility
		axiMem[axiOffset/4+uint32(i)] = uint32(key[i])
	}

	return nil
}

// LoadEncryptedCode loads encrypted code into FPGA memory via AXI
func LoadEncryptedCode(encryptedCode []byte, iv []byte, axiOffset uint32, mappedMem []byte) error {
	// Ensure code and IV fit into mapped memory
	if len(encryptedCode)+len(iv) > len(mappedMem) {
		return fmt.Errorf("encrypted code size exceeds mapped memory")
	}

	// Load IV to memory (you can set a specific memory region for IV if needed)
	for i := 0; i < len(iv); i++ {
		mappedMem[axiOffset+uint32(i)] = iv[i]
	}

	// Load encrypted code to FPGA memory
	for i := 0; i < len(encryptedCode); i++ {
		mappedMem[axiOffset+uint32(i+len(iv))] = encryptedCode[i]
	}

	return nil
}

// ExecuteDecryptedCode sends a command to the FPGA to decrypt and execute code
func ExecuteDecryptedCode(mappedMem []byte, commandOffset uint32) error {
	// Write to the control register (this address may vary based on your FPGA design)
	mappedMem[commandOffset] = 1 // Set '1' to start decryption and execution

	// Optionally wait for completion (polling or interrupts can be used here)
	// For polling example:
	for {
		if mappedMem[commandOffset] == 0 {
			break
		}
	}

	fmt.Println("Decryption and execution completed on FPGA")
	return nil
}
