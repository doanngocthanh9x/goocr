#!/bin/bash

# Script chạy Python gRPC server

set -e

echo "Starting Python gRPC Server..."
echo "Server will listen on $(hostname -I):50051"

cd python_server
python server.py
