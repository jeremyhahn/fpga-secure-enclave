# Makefile for compiling and simulating Verilog source files, and for building and testing Golang code.

# Variables for Verilog Compilation
VLOG_COMPILER = iverilog
VLOG_SIMULATOR = vvp
VLOG_WAVEFORM_VIEWER = gtkwave
VLOG_TARGET = build/enclave_sim
VLOG_SOURCES = $(wildcard verilog/**/*.v)
VLOG_TOP_MODULE = rocket_chip_enclave  # Specify your top-level Verilog module
VLOG_WAVEFORM_OUTPUT = build/waveform.vcd

# Variables for Golang Compilation and Testing
GO_COMPILER = go
GO_BUILD_DIR = build
GO_MAIN_FILE = cmd/enclave/main.go
GO_TEST_DIR = tests

# Default target: Compile, simulate Verilog, and build Golang code
all: compile_vlog simulate_vlog build_go

# Verilog targets
compile_vlog:
	@mkdir -p $(GO_BUILD_DIR)
	$(VLOG_COMPILER) -o $(VLOG_TARGET) $(VLOG_SOURCES) -s $(VLOG_TOP_MODULE)
	@echo "Verilog compilation complete."

simulate_vlog: compile_vlog
	$(VLOG_SIMULATOR) $(VLOG_TARGET)
	@echo "Verilog simulation complete."

waveform_vlog: compile_vlog
	$(VLOG_SIMULATOR) $(VLOG_TARGET) -lxt2
	$(VLOG_WAVEFORM_VIEWER) $(VLOG_WAVEFORM_OUTPUT)
	@echo "Waveform generated and loaded in GTKWave."

clean_vlog:
	rm -rf $(GO_BUILD_DIR)
	@echo "Verilog build directory cleaned."

# Golang targets
build_go:
	@echo "Building Golang code..."
	$(GO_COMPILER) build -o $(GO_BUILD_DIR)/enclave $(GO_MAIN_FILE)
	@echo "Golang build complete."

test_go:
	@echo "Running Golang unit tests..."
	$(GO_COMPILER) test $(GO_TEST_DIR) -v
	@echo "Golang unit tests complete."

clean_go:
	@echo "Cleaning Golang build artifacts..."
	rm -rf $(GO_BUILD_DIR)/enclave
	@echo "Golang build artifacts cleaned."

# Clean all build artifacts (Verilog and Golang)
clean: clean_vlog clean_go

# Help target
help:
	@echo "Makefile Targets:"
	@echo "  all           - Compile and simulate Verilog, and build Golang code"
	@echo "  compile_vlog  - Compile Verilog source files"
	@echo "  simulate_vlog - Simulate compiled Verilog code"
	@echo "  waveform_vlog - Generate waveform and view with GTKWave"
	@echo "  clean_vlog    - Clean Verilog build directory"
	@echo "  build_go      - Build Golang code"
	@echo "  test_go       - Run Golang unit tests"
	@echo "  clean_go      - Clean Golang build artifacts"
	@echo "  clean         - Clean all build artifacts (Verilog and Golang)"
	@echo "  help          - Show this help message"
