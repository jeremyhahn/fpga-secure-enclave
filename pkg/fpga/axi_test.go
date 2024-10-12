package fpga

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadKeyToFPGA(t *testing.T) {
	// Simulated FPGA memory
	mappedMem := make([]byte, 1024)

	// Example key (256 bits)
	key := []byte("ThisIsA32ByteKeyForAES256Encryption!")
	axiOffset := uint32(0)

	err := LoadKeyToFPGA(key, axiOffset, mappedMem)
	assert.Nil(t, err)
	assert.Equal(t, key, mappedMem[axiOffset:axiOffset+uint32(len(key))])
}

func TestLoadEncryptedCode(t *testing.T) {
	// Simulated FPGA memory
	mappedMem := make([]byte, 1024)

	// Example encrypted code and IV
	encryptedCode := []byte("Encrypted code data")
	iv := []byte("InitializationVec")
	axiOffset := uint32(0)

	err := LoadEncryptedCode(encryptedCode, iv, axiOffset, mappedMem)
	assert.Nil(t, err)
	assert.Equal(t, iv, mappedMem[axiOffset:axiOffset+uint32(len(iv))])
	assert.Equal(t, encryptedCode, mappedMem[axiOffset+uint32(len(iv)):axiOffset+uint32(len(iv)+len(encryptedCode))])
}

func TestExecuteDecryptedCode(t *testing.T) {
	// Simulated FPGA memory
	mappedMem := make([]byte, 1024)

	// Simulated commandOffset
	commandOffset := uint32(500)

	// Test executing code on FPGA
	err := ExecuteDecryptedCode(mappedMem, commandOffset)
	assert.Nil(t, err)
	assert.Equal(t, byte(0), mappedMem[commandOffset]) // Ensure command is reset after execution
}
