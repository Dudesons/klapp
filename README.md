# Klapp

## Flip model

Flip values available:
 * boolean (with the activated field in the payload)
 * int
 * string
 * slice of string
 * slice of int

Fields:
 * activated (boolean) -- *REQUIRED*: Determine if the flip is activated
 * name (string) -- *REQUIRED*: 
 * description (string) -- highly recommended: 
 * type (integer) -- *REQUIRED for slice values*: Indicate if flip values are type string (0) or int (1)
 * content -- *REQUIRED (except for boolean flags)*: This is the content of a feature flag, the content can be an integer, string, int slice, string slice 


protobuff with grpc: `docker run --rm -v $(pwd):$(pwd) -w $(pwd) znly/protoc --go_out=plugins=grpc:. -I. ./src/github.com/tinyclues/guay/proto/api.proto`
setup your host/vm for running elasticsearch with: `sudo sysctl -w vm.max_map_count=262144`

```
docker run --rm -v $(pwd):$(pwd) -w $(pwd) znly/protoc -I$GOPATH/src \
                                                       -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
                                                       --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
                                                       -I. ./proto/api.proto
```

docker run --rm -v $(pwd):$(pwd) -w $(pwd) znly/protoc -I$GOPATH/src \
                                                       --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
                                                       -I. ./api.proto
                             
docker run --rm -v $(pwd):$(pwd) -w $(pwd) znly/protoc -I$GOPATH/src --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. -I. ./pb/api.proto                                                       
docker run --rm -v $(pwd):$(pwd) -w $(pwd) znly/protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  --grpc-gateway_out=logtostderr=true:. ./pb/api.proto
docker run --rm -v $(pwd):$(pwd) -w $(pwd) znly/protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  --swagger_out=logtostderr=true:. ./pb/api.proto

## GOOD
sudo docker run --rm -v $(pwd):$(pwd) -w $(pwd) znly/protoc -I/usr/local/include -I. -I${GOPATH}/src --go_out=plugins=grpc:. ./pb/api.proto


sudo docker run --rm -v $(pwd):$(pwd) -w $(pwd) znly/protoc -I/usr/local/include -I. \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. -I ./pb/api.proto
sudo docker run --rm -v $(pwd):$(pwd) -w $(pwd) znly/protoc -I/usr/local/include -I. \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:. -I ./pb/api.proto                                        