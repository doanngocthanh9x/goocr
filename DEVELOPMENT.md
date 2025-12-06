# Development Guide - GoOCR

Hướng dẫn chi tiết cho các developers muốn mở rộng hoặc maintain project.

## 📦 Dependency Management

### Python Dependencies

File: `requirements.txt`

```
grpcio==1.59.0              # gRPC framework
grpcio-tools==1.59.0        # gRPC code generation tools
protobuf==4.24.4            # Protocol buffers library
```

**Cập nhật dependencies:**

```bash
# Add new package
pip install package_name
pip freeze > requirements.txt

# Update existing package
pip install --upgrade package_name
pip freeze > requirements.txt

# Install all from requirements
pip install -r requirements.txt
```

### Go Dependencies

Files: `go.mod`, `go.work`, `datatransfer/go.mod`, `common/go.mod`

**Go Workspace Structure:**

```
go.work (root)
  ├── . (main module: github.com/doanngocthanh9x/goocr)
  ├── ./common (github.com/doanngocthanh9x/goocr/gen/common)
  └── ./datatransfer (github.com/doanngocthanh9x/goocr/gen/datatransfer)
```

**Cập nhật Go dependencies:**

```bash
# Download all dependencies
go mod download

# Tidy up go.mod (remove unused, add missing)
go mod tidy

# View current dependencies
go list -m all

# Update specific package
go get -u github.com/package@latest

# Check for updates
go list -u -m all
```

**Module Paths (không được thay đổi):**

```
github.com/doanngocthanh9x/goocr              # Main module
github.com/doanngocthanh9x/goocr/gen/common   # Generated common (→ ./common)
github.com/doanngocthanh9x/goocr/gen/datatransfer  # Generated datatransfer (→ ./datatransfer)
```

## 🔧 Build & Development Workflow

### Generate gRPC Code

Proto file locations: `proto/`

Generated code locations:
- **Python**: `python_server/generated/`
- **Go**: `common/`, `datatransfer/`

```bash
# Generate
make proto-gen
# hoặc: ./scripts/gen_proto.sh

# Clean generated code
make clean
```

### Python Development

```bash
# Run server
make server
# hoặc: cd python_server && python server.py

# Check Python version
python --version

# Install development tools
pip install pylint black pytest
```

**Add new handler:**

1. Create file: `python_server/handlers/new_handler.py`
2. Implement handler class
3. Import in `python_server/server.py`
4. Add servicer method in `DataTransferServicer` class

### Go Development

```bash
# Build client binary
make build

# Run client
make client
# hoặc: go run ./cmd/client/...

# Check Go version
go version

# View workspace status
go work use -r .

# Build for Linux
CGO_ENABLED=0 GOOS=linux go build -o bin/client-linux ./cmd/client/...

# Build for Windows
CGO_ENABLED=0 GOOS=windows go build -o bin/client.exe ./cmd/client/...
```

**Add new client type:**

1. Add proto definition in `proto/data_transfer.proto`
2. Run `make proto-gen`
3. Create client file: `cmd/client/new_client.go`
4. Implement client structure and methods
5. Call in `cmd/client/main.go`

## 🧪 Testing

### Python Tests

```bash
# Install pytest
pip install pytest

# Create test file: python_server/tests/test_handlers.py
# Run tests
pytest python_server/tests/ -v
```

Example test:

```python
# python_server/tests/test_file_handler.py
from python_server.handlers.file_handler import FileHandler

def test_process_file():
    test_data = b"test content"
    result = FileHandler.process_file(test_data, "test.txt")
    
    assert result["filename"] == "test.txt"
    assert result["size"] == 12
    assert result["status"] == "processed"
```

### Go Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...
```

Example test:

```go
// cmd/client/client_test.go
package main

import "testing"

func TestNewFileClient(t *testing.T) {
    // Test implementation
}
```

## 📝 Protocol Buffer Changes

### Adding New RPC Method

**Step 1:** Update proto file

```protobuf
// proto/data_transfer.proto
service DataTransferService {
  rpc ProcessFile(stream FileChunk) returns (ProcessFileResponse) {}
  rpc ProcessJson(JsonRequest) returns (JsonResponse) {}
  rpc ProcessImage(ImageRequest) returns (ImageResponse) {}  // NEW
}

message ImageRequest {
  bytes image_data = 1;
  string filename = 2;
}

message ImageResponse {
  ProcessingStatus status = 1;
  string result = 2;
}
```

**Step 2:** Generate code

```bash
make proto-gen
```

**Step 3:** Implement Python handler

```python
# python_server/handlers/image_handler.py
class ImageHandler:
    @staticmethod
    def process_image(image_data: bytes, filename: str) -> dict:
        # Implementation
        return {"result": "processed"}
```

**Step 4:** Add Python servicer method

```python
# python_server/server.py
def ProcessImage(self, request, context):
    try:
        result = ImageHandler.process_image(request.image_data, request.filename)
        response = data_transfer_pb2.ImageResponse()
        response.status.status = common_pb2.ProcessingStatus.SUCCESS
        response.result = str(result)
        return response
    except Exception as e:
        response = data_transfer_pb2.ImageResponse()
        response.status.status = common_pb2.ProcessingStatus.ERROR
        response.status.message = str(e)
        return response
```

**Step 5:** Implement Go client

```go
// cmd/client/client.go
type ImageClient struct {
    conn *grpc.ClientConn
}

func NewImageClient(conn *grpc.ClientConn) *ImageClient {
    return &ImageClient{conn: conn}
}

func (ic *ImageClient) SendImage(imagePath string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
    defer cancel()

    client := pb.NewDataTransferServiceClient(ic.conn)
    
    // Read image file
    imageData, err := os.ReadFile(imagePath)
    if err != nil {
        return "", err
    }

    request := &pb.ImageRequest{
        ImageData: imageData,
        Filename:  filepath.Base(imagePath),
    }

    response, err := client.ProcessImage(ctx, request)
    if err != nil {
        return "", err
    }

    return response.Result, nil
}
```

**Step 6:** Update Go client main

```go
// cmd/client/main.go
func main() {
    // ... existing code ...
    
    // Test 3: Image processing
    imageClient := NewImageClient(conn)
    result, err := imageClient.SendImage("./test_image.jpg")
    // ... handle result ...
}
```

## 📊 Project Statistics

```bash
# Lines of code
find . -name "*.py" -o -name "*.go" -o -name "*.proto" | xargs wc -l

# File count
find . -name "*.py" -o -name "*.go" | wc -l

# Git status
git status
git log --oneline
```

## 🚀 Deployment

### Docker Setup (Optional)

```dockerfile
# Dockerfile.python
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY python_server ./python_server
COPY proto ./proto
COPY scripts ./scripts
RUN chmod +x scripts/gen_proto.sh && ./scripts/gen_proto.sh
EXPOSE 50051
CMD ["python", "python_server/server.py"]
```

### Running in Docker

```bash
# Build
docker build -f Dockerfile.python -t goocr-server .

# Run
docker run -p 50051:50051 goocr-server
```

## 🐛 Troubleshooting

### protoc not found

```bash
# Ubuntu/Debian
sudo apt install protobuf-compiler

# macOS
brew install protobuf

# Check version
protoc --version
```

### Go module issues

```bash
# Clean cache
go clean -modcache

# Tidy modules
go mod tidy

# Verify modules
go mod verify
```

### Python import errors

```bash
# Reinstall
pip install --force-reinstall -r requirements.txt

# Check installation
python -c "import grpc; print(grpc.__version__)"
```

### Server connection refused

```bash
# Check if server is running
lsof -i :50051

# Check firewall
sudo ufw status

# Kill existing process if needed
pkill -f "python server.py"
```

## 📚 References

- [gRPC Documentation](https://grpc.io/docs/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [Go gRPC](https://grpc.io/docs/languages/go/)
- [Python gRPC](https://grpc.io/docs/languages/python/)
- [Go Modules](https://go.dev/blog/using-go-modules)
