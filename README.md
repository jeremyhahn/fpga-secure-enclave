# FPGA Secure Enclave

## Project Overview

This project implements an experimental secure enclave on an FPGA platform, supporting cryptographic operations such as signing, encryption, and decryption. It integrates the Rocket Chip core, Shamir Secret Sharing for key splitting and partial signing, and includes a Golang client that communicates with the FPGA via the AXI interface.

### Features

- **Cryptographic Algorithms**: Supports RSA, ECDSA, Ed25519, and AES-256.
- **Shamir Secret Sharing**: Key splitting via Shamir Secret Sharing with threshold signing support.
- **Memory Encryption**: AES-256 encryption for instruction and memory protection.
- **Tamper Resistance**: Tamper-resistant key storage on the FPGA.
- **Secure World**: Execute code in a "secure world" similar to ARM TrustZone or Intel SGX
- **Golang Client**: Manages secure communication, key & code loading, and cryptographic operations.


# Modules Description

## Verilog Modules

- **rocket_chip/rocket_chip_enclave.v**: Integrates the Rocket Chip core with the secure enclave for encrypted code execution.
- **aes/aes256_ctr.v**: AES-256 encryption and decryption in CTR mode.
- **rsa/rsa_signing_core.v**: RSA signing core with support for both full and partial-key signing.
- **ecdsa/ecdsa_signing_core.v**: ECDSA signing core with support for both full and partial-key signing.
- **ed25519/ed25519_signing_core.v**: Ed25519 signing core with support for both full and partial-key signing.
- **tamper_detection/key_storage_with_tamper.v**: Tamper-resistant storage for cryptographic keys.

## Golang Modules

- **enclave/aes.go**: Manages AES-256 key initialization, encryption, and decryption functions.
- **enclave/rsa.go**: Manages RSA key initialization, Shamir Secret Sharing for key splitting, and signing functions.
- **enclave/ecdsa.go**: Manages ECDSA key initialization, Shamir Secret Sharing for key splitting, and signing functions.
- **enclave/ed25519.go**: Manages Ed25519 key initialization, Shamir Secret Sharing for key splitting, and signing functions.
- **enclave/enclave.go**: Handles enclave initialization and secure key loading.
- **fpga/axi.go**: Handles AXI communication between the Golang client and the FPGA.
- **fpga/memory.go**: Manages memory for FPGA key storage.

# Build

    # Install dependencies
    sudo apt-get install iverilog gtkwave golang scala sbt verilator

### Rocket Chip Core

Once you have cloned the repository, navigate to the rocket-chip directory and build the project using SBT:

    git clone https://github.com/chipsalliance/rocket-chip.git
    cd rocket-chip
    git submodule update --init
    cd emulator
    make


### Verilog Compilation and Simulation

    # Verilog Compilation and Simulation
    make compile_vlog

    # Run Verilog simulation:
    make simulate_vlog

    # View waveform in GTKWave
    make waveform_vlog

    
### Golang Build and Run

    # Build the Golang client
    make build_go
    ./build/enclave


# Loading the Secure Enclave onto the FPGA

1. Map FPGA memory via /dev/mem. This is handled by the Golang client.
2. Load cryptographic keys (AES, RSA, ECDSA, Ed25519) onto the FPGA using the AXI interface.
3. Run the enclave after key loading. The enclave will then be ready for secure cryptographic operations.


# Performing Signing Operations

### AES-256 Encryption
```go
plaintext := []byte("Test data for AES encryption.")
ciphertext, err := enclave.AESEncrypt(plaintext, keyStore)
if err != nil {
    log.Fatalf("AES encryption failed: %v", err)
}
fmt.Printf("AES Ciphertext: %x\n", ciphertext)
```

### AES-256 Decryption
```go
decryptedText, err := enclave.AESDecrypt(ciphertext, keyStore)
if err != nil {
    log.Fatalf("AES decryption failed: %v", err)
}
fmt.Printf("Decrypted Text: %s\n", string(decryptedText))
```

### RSA Full Signing
```go
message := []byte("Test message for signing.")
signature, err := enclave.RSASign(message, keyStore)
if err != nil {
    log.Fatalf("RSA full signing failed: %v", err)
}
fmt.Printf("RSA Full Signature: %x\n", signature)
```

### RSA Partial Signing

```go
partialSignature, err := enclave.RSAPartialSign(message, keyStore)
if err != nil {
    log.Fatalf("RSA partial signing failed: %v", err)
}
fmt.Printf("RSA Partial Signature: %x\n", partialSignature)
```

### ECDSA Full Signing
```go
ecdsaSignature, err := enclave.ECDSASign(message, keyStore)
if err != nil {
    log.Fatalf("ECDSA full signing failed: %v", err)
}
fmt.Printf("ECDSA Full Signature: %x\n", ecdsaSignature)
```

### ECDSA Partial Signing
```go
ecdsaPartialSignature, err := enclave.ECDSAPartialSign(message, keyStore)
if err != nil {
    log.Fatalf("ECDSA partial signing failed: %v", err)
}
fmt.Printf("ECDSA Partial Signature: %x\n", ecdsaPartialSignature)
```

### Ed25519 Full Signing
```go
ed25519Signature, err := enclave.Ed25519Sign(message, keyStore)
if err != nil {
    log.Fatalf("Ed25519 full signing failed: %v", err)
}
fmt.Printf("Ed25519 Full Signature: %x\n", ed25519Signature)
```

### Ed25519 Partial Signing
```go
ed25519PartialSignature, err := enclave.Ed25519PartialSign(message, keyStore)
if err != nil {
    log.Fatalf("Ed25519 partial signing failed: %v", err)
}
fmt.Printf("Ed25519 Partial Signature: %x\n", ed25519PartialSignature)
```

# Unit Testing

    make test_go

# Dependencies

This project makes use of the following IP cores:
- https://github.com/secworks/aes
- https://github.com/chipsalliance/rocket-chip.git
