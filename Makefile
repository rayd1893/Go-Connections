lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint
	@golangci-lint run --config .golangci.yaml
.PHONY: lint

protoc:
	@go install golang.org/x/tools/cmd/goimports google.golang.org/grpc/cmd/protoc-gen-go-grpc google.golang.org/protobuf/cmd/protoc-gen-go
	@protoc --proto_path=. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./pkg/pb/*.proto
	@goimports -w pkg/pb
.PHONY: protoc
