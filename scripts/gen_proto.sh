#!/bin/bash

# Script sinh code từ proto files

set -e

PROTO_DIR="./proto"
PYTHON_GEN_DIR="./python_server/generated"
GO_GEN_DIR="./go_client/generated"

echo "Generating Python gRPC code..."
python -m grpc_tools.protoc \
    -I${PROTO_DIR} \
    --python_out=${PYTHON_GEN_DIR} \
    --grpc_python_out=${PYTHON_GEN_DIR} \
    ${PROTO_DIR}/*.proto

echo "Generating Go gRPC code..."
protoc \
    -I${PROTO_DIR} \
    --go_out=${GO_GEN_DIR} \
    --go-grpc_out=${GO_GEN_DIR} \
    --go_opt=module=github.com/doanngocthanh9x/goocr \
    --go-grpc_opt=module=github.com/doanngocthanh9x/goocr \
    ${PROTO_DIR}/*.proto

echo "✓ Proto code generation completed successfully!"
echo "  - Python stubs: ${PYTHON_GEN_DIR}"
echo "  - Go stubs: ${GO_GEN_DIR}"
