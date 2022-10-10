proto_def = server/kong_plugin_protocol/pluginsocket.pb.go

.PHONY: lint
lint:
	golangci-lint run --exclude composites

.PHONY: dep
dep: $(proto_def)
	go get -v
	go mod tidy

.PHONY: test
test: dep
	go test -v -race ./...

pluginsocket.proto:
	wget https://raw.githubusercontent.com/Kong/kong/master/kong/include/kong/pluginsocket.proto

$(proto_def): pluginsocket.proto
	mkdir -p server/kong_plugin_protocol
	protoc -I . pluginsocket.proto --go_out=./server/kong_plugin_protocol/ --go_opt=paths=source_relative
