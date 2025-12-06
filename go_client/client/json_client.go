package main

import (
	"context"
	"log"
	"time"

	pb "github.com/doanngocthanh9x/goocr/gen/datatransfer"

	"google.golang.org/grpc"
)

const contextTimeout = 30 * time.Second

// JsonClient - gRPC client cho xử lý JSON
type JsonClient struct {
	conn *grpc.ClientConn
}

// NewJsonClient - Tạo mới JsonClient
func NewJsonClient(conn *grpc.ClientConn) *JsonClient {
	return &JsonClient{
		conn: conn,
	}
}

// SendJson - Gửi JSON tới server
func (jc *JsonClient) SendJson(jsonData string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	client := pb.NewDataTransferServiceClient(jc.conn)

	request := &pb.JsonRequest{
		JsonData: jsonData,
	}

	log.Printf("Sending JSON request: %s", jsonData)

	response, err := client.ProcessJson(ctx, request)
	if err != nil {
		log.Printf("Error sending JSON: %v", err)
		return "", err
	}

	log.Printf("Response status: %v, message: %s", response.Status.Status, response.Status.Message)
	return response.Result, nil
}
