package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// LoadConfig - Load config từ environment
type Config struct {
	ServerHost string
	ServerPort int
}

func LoadConfig() Config {
	host := os.Getenv("GRPC_SERVER_HOST")
	if host == "" {
		host = "localhost"
	}

	portStr := os.Getenv("GRPC_SERVER_PORT")
	port := 50051
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	return Config{
		ServerHost: host,
		ServerPort: port,
	}
}

func main() {
	// Load config
	config := LoadConfig()
	serverAddr := fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
	log.Printf("Connecting to server: %s", serverAddr)

	// Connect to server
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to gRPC server")

	// Test 1: Gửi JSON
	log.Println("\n=== Test 1: Processing JSON ===")
	jsonClient := NewJsonClient(conn)
	jsonData := `{"name": "John", "age": 30, "email": "john@example.com"}`
	result, err := jsonClient.SendJson(jsonData)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Result: %s", result)
	}

	// Test 2: Gửi file
	log.Println("\n=== Test 2: Processing File ===")
	fileClient := NewFileClient(conn)
	
	// Tạo file test
	testFile := "/tmp/test.txt"
	if err := createTestFile(testFile); err != nil {
		log.Fatalf("Error creating test file: %v", err)
	}
	
	result, err = fileClient.SendFile(testFile)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Result: %s", result)
	}
}

// createTestFile - Tạo file test
func createTestFile(filePath string) error {
	testData := []byte("Hello, this is a test file for gRPC file transfer.\nLine 2\nLine 3\n")
	err := os.WriteFile(filePath, testData, 0644)
	if err != nil {
		return err
	}
	log.Printf("Test file created: %s", filePath)
	return nil
}
