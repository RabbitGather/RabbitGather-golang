
.PHONY: protof
protof:
	protoc --proto_path=./proto/proto --go_out=./proto/pb --go-grpc_out=./proto/pb ./proto/proto/*