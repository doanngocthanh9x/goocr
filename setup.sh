#!/bin/bash

# Setup script for GoOCR project

set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║   GoOCR - gRPC File & JSON Transfer System Setup          ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Check prerequisites
echo "📋 Checking prerequisites..."

# Check Python
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3 is not installed. Please install Python 3.8+"
    exit 1
fi
echo "✓ Python $(python3 --version | cut -d' ' -f2) found"

# Check Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.19+"
    exit 1
fi
echo "✓ Go $(go version | cut -d' ' -f3) found"

# Install protoc if not exists
if ! command -v protoc &> /dev/null; then
    echo ""
    echo "📥 Installing protoc..."
    sudo apt update
    sudo apt install -y protobuf-compiler
fi
echo "✓ protoc $(protoc --version | cut -d' ' -f2) found"

# Install Go protoc plugins
echo ""
echo "📥 Installing Go protoc plugins..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
echo "✓ Go protoc plugins installed"

# Install Python dependencies
echo ""
echo "📥 Installing Python dependencies..."
pip install -r requirements.txt
echo "✓ Python dependencies installed"

# Download Go modules
echo ""
echo "📥 Downloading Go modules..."
go mod download
echo "✓ Go modules downloaded"

# Generate proto code
echo ""
echo "🔧 Generating gRPC code from proto files..."
chmod +x scripts/gen_proto.sh
./scripts/gen_proto.sh
echo "✓ gRPC code generated"

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║              ✓ Setup Complete!                             ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "Next steps:"
echo "  1. Terminal 1: make server"
echo "  2. Terminal 2: make client"
echo ""
echo "Or use: make help"
