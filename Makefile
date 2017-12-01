help:
	@echo "proto: Builds the protobuf"

build:
	docker run --rm -v $(pwd):$(pwd) -v $GOPATH/src/github.com/grpc-ecosystem:$GOPATH/src/github.com/grpc-ecosystem -w $(pwd) znly/protoc -I$GOPATH/src --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. -I. ./pb/api.proto
	docker run --rm -v $(pwd):$(pwd) -v $GOPATH/src/github.com/grpc-ecosystem:$GOPATH/src/github.com/grpc-ecosystem -w $(pwd) znly/protoc -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  --grpc-gateway_out=logtostderr=true:. -I. ./pb/api.proto
	docker run --rm -v $(pwd):$(pwd) -v $GOPATH/src/github.com/grpc-ecosystem:$GOPATH/src/github.com/grpc-ecosystem -w $(pwd) znly/protoc -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  --swagger_out=logtostderr=true:. -I. ./pb/api.proto
