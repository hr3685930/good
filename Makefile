proto:
	protoc --go_out=plugins=grpc:. api/proto/v1/*/*.proto