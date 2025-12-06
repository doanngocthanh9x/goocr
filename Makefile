.PHONY: help setup proto-gen server client build test clean

help:
	@echo "GoOCR - gRPC File & JSON Transfer System"
	@echo ""
	@echo "Available commands:"
	@echo "  make setup        - Setup environment (install dependencies)"
	@echo "  make proto-gen    - Generate gRPC code from proto files"
	@echo "  make server       - Run Python gRPC server"
	@echo "  make client       - Run Go client (requires server running)"
	@echo "  make build        - Build Go client binary"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean generated files and binaries"
	@echo ""

# Setup environment
setup:
	@echo "Setting up environment..."
	@echo "Installing Python dependencies..."
	pip install -r requirements.txt
	@echo "Downloading Go modules..."
	go mod download
	@echo "✓ Environment setup complete"

# Generate gRPC code from proto
proto-gen:
	@echo "Generating gRPC code from proto files..."
	./scripts/gen_proto.sh
	@echo "✓ Proto code generation complete"

# Run Python gRPC server
server:
	@echo "Starting Python gRPC Server..."
	./scripts/run_server.sh

# Run Go client
client:
	@echo "Running Go client..."
	go run ./cmd/client/...

# Build Go client
build:
	@echo "Building Go client..."
	mkdir -p bin
	go build -o bin/client ./cmd/client/...
	@echo "✓ Build complete: bin/client"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f python_server/generated/*pb2*.py
	rm -f common/*.pb.go datatransfer/*pb.go datatransfer/*grpc.pb.go
	@echo "✓ Cleanup complete"

# Install protocol buffers compiler
install-protoc:
	@echo "Installing protoc..."
	@which protoc > /dev/null || (sudo apt update && sudo apt install -y protobuf-compiler)
	@echo "✓ protoc installed"

# Install Go protoc plugins
install-go-protoc-plugins:
	@echo "Installing Go protoc plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✓ Go protoc plugins installed"

# Full setup (one command)
fullsetup: install-protoc install-go-protoc-plugins setup proto-gen
	@echo "✓ Full setup complete! Run 'make server' and 'make client' in separate terminals"
