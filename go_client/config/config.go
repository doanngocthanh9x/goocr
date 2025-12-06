package main

import (
	"os"
	"strconv"
)

// Config - Cấu hình Go client
type Config struct {
	ServerHost string
	ServerPort int
}

// LoadConfig - Load config từ environment
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
