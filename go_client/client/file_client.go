package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/doanngocthanh9x/goocr/gen/datatransfer"

	"google.golang.org/grpc"
)

// FileClient - gRPC client cho xử lý file
type FileClient struct {
	client pb.DataTransferService_ProcessFileClient
	conn   *grpc.ClientConn
}

// NewFileClient - Tạo mới FileClient
func NewFileClient(conn *grpc.ClientConn) *FileClient {
	return &FileClient{
		conn: conn,
	}
}

// SendFile - Gửi file tới server
func (fc *FileClient) SendFile(filePath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	client := pb.NewDataTransferServiceClient(fc.conn)

	// Mở file
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return "", err
	}
	defer file.Close()

	// Tạo stream client
	stream, err := client.ProcessFile(ctx)
	if err != nil {
		log.Printf("Error creating stream: %v", err)
		return "", err
	}

	// Đọc file và gửi chunks
	filename := file.Name()
	chunkSize := 1024 * 1024 // 1MB per chunk
	chunkIndex := 0
	totalChunks := 0

	// Đầu tiên đếm số chunks
	fileInfo, err := file.Stat()
	if err == nil {
		totalChunks = int((fileInfo.Size() + int64(chunkSize) - 1) / int64(chunkSize))
	}

	// Reset file pointer
	file.Seek(0, 0)

	// Gửi chunks
	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			log.Printf("Error reading file: %v", err)
			return "", err
		}

		if n == 0 {
			break
		}

		chunk := &pb.FileChunk{
			Filename:   filename,
			Data:       buffer[:n],
			ChunkIndex: int32(chunkIndex),
			TotalChunks: int32(totalChunks),
		}

		if err := stream.Send(chunk); err != nil {
			log.Printf("Error sending chunk: %v", err)
			return "", err
		}

		log.Printf("Sent chunk %d/%d", chunkIndex+1, totalChunks)
		chunkIndex++

		if err == io.EOF {
			break
		}
	}

	// Nhận response
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("Error closing stream: %v", err)
		return "", err
	}

	log.Printf("Response status: %v, message: %s", response.Status.Status, response.Status.Message)
	return response.Result, nil
}
