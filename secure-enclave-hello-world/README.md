# Secure Enclave - Hello World

Example Hello World Program for the Secure Enclave. This is a simple "Hello World" program written in C that can be compiled and run on the Rocket Chip secure enclave core. This program prints "Hello, Secure Enclave!" to the output.


# Instructions to Compile and Load onto the FPGA

1. Install RISC-V Toolchain

    # The Rocket Chip core uses the RISC-V architecture, so you will need
    # the RISC-V toolchain to compile the program. Install the RISC-V
    # toolchain using the following commands:
    git clone https://github.com/riscv/riscv-gnu-toolchain
    cd riscv-gnu-toolchain
    ./configure --prefix=/opt/riscv
    make

    # Ensure the path to the toolchain binaries is in your environment
    export PATH=/opt/riscv/bin:$PATH

2. Compile the Hello World Program

    # Compile the program
    riscv64-unknown-elf-gcc -o hello hello.c
    riscv64-unknown-elf-objcopy -O binary hello hello.bin

3. Encrypt the Hello World Binary

    Use the AES-256 encryption from the Golang code to encrypt hello.bin before loading it onto the FPGA

```go
package main

import (
	"fmt"
	"io/ioutil"
	"github.com/jeremyhahn/fpga-secure-enclave/enclave"
)

func main() {
	key := []byte("ThisIsA32ByteKeyForAES256Encryption!")
	plaintext, err := ioutil.ReadFile("hello.bin")
	if err != nil {
		fmt.Println("Error reading binary file:", err)
		return
	}

	encryptedCode, iv, err := enclave.EncryptCodeAES(plaintext, key)
	if err != nil {
		fmt.Println("Error encrypting code:", err)
		return
	}

	// Save encrypted code and IV for loading into the FPGA
	ioutil.WriteFile("encrypted_hello.bin", encryptedCode, 0644)
	ioutil.WriteFile("iv.bin", iv, 0644)

	fmt.Println("Code encrypted successfully!")
}
```

4. Load the Encrypted Program into the FPGA
```go
package main

import (
	"fmt"
	"io/ioutil"
	"syscall"
	"github.com/jeremyhahn/fpga-secure-enclave/fpga"
)

const axiOffset = 0x2000
const commandOffset = 0x3000

func main() {
	// Open the encrypted code and IV
	encryptedCode, _ := ioutil.ReadFile("encrypted_hello.bin")
	iv, _ := ioutil.ReadFile("iv.bin")

	// Memory map the FPGA
	fd, err := syscall.Open("/dev/mem", syscall.O_RDWR|syscall.O_SYNC, 0)
	if err != nil {
		fmt.Printf("Failed to open /dev/mem: %v\n", err)
		return
	}
	defer syscall.Close(fd)

	pageSize := syscall.Getpagesize()
	mappedMem, err := syscall.Mmap(fd, 0, pageSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		fmt.Printf("Failed to mmap FPGA memory: %v\n", err)
		return
	}
	defer syscall.Munmap(mappedMem)

	// Load encrypted hello program into the FPGA
	err = fpga.LoadEncryptedCode(encryptedCode, iv, axiOffset, mappedMem)
	if err != nil {
		fmt.Printf("Failed to load encrypted code into FPGA: %v\n", err)
		return
	}

	// Execute the loaded program on the secure enclave
	err = fpga.ExecuteDecryptedCode(mappedMem, commandOffset)
	if err != nil {
		fmt.Printf("Execution failed: %v\n", err)
	}
}
```

5. Run the Program on the Secure Enclave

Once the encrypted code is loaded and the execution command is issued, the FPGA will decrypt the binary and execute it. You should see the "Hello, Secure Enclave!" output.