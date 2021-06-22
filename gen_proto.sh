#/bin/bash
mkdir -p ./app/hcproto
protoc --go_out=./app/hcproto \
    --go-grpc_out=./app/hcproto \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative healthcheck.proto
