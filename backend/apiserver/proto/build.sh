#!/bin/bash

find . -name '*.pb.go' -exec rm -rf {} \; -o -name '*.swagger.json' -exec rm -rf {} \; -o -name '*.gw.go' -exec rm -rf {} \;

export GO111MODULE="on"


GW_ORIGIN_VERSION=v1.14.4

# plugin for go porotbuf
go get google.golang.org/protobuf/cmd/protoc-gen-go
# plugin for go grpc
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
# plugin for go grpc-gateway
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway 
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 
go get github.com/grpc-ecosystem/grpc-gateway@$GW_ORIGIN_VERSION


GOPATH=$(go env GOPATH)
ARGS="\
--proto_path=. \
--proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$GW_ORIGIN_VERSION \
--proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$GW_ORIGIN_VERSION/third_party/googleapis \
--go_out=. \
--go_opt=paths=source_relative \
--go-grpc_out=. \
--go-grpc_opt=paths=source_relative \
--grpc-gateway_opt=logtostderr=true \
--grpc-gateway_out=. \
--grpc-gateway_opt=paths=source_relative \
--grpc-gateway_opt=repeated_path_param_separator=ssv"

for dir in $(find . -name '*.proto' | xargs -I{} dirname {} | sort | uniq); do
   echo "building $dir/*.proto"
   protoc $ARGS $dir/*.proto
done