#!/bin/bash -e

go install google.golang.org/protobuf/cmd/protoc-gen-go
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

curl -o $TMP_DIR/pluginsocket.proto https://raw.githubusercontent.com/Kong/kong/master/kong/pluginsocket.proto

protoc --go_out=$TMP_DIR --proto_path=$TMP_DIR pluginsocket.proto

cp $TMP_DIR/kong_plugin_protocol/pluginsocket.pb.go server/kong_plugin_protocol/pluginsocket.pb.go