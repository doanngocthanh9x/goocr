module github.com/doanngocthanh9x/goocr/gen/datatransfer

go 1.21

require (
	github.com/doanngocthanh9x/goocr/gen/common v0.0.0
	google.golang.org/grpc v1.77.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250106144421-5f2ede2c6b0b // indirect
)

replace github.com/doanngocthanh9x/goocr/gen/common => ../common
