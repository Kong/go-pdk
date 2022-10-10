proto_def = ./server/kong_plugin_protocol/pluginsocket.proto
proto_def_compiled = ./server/kong_plugin_protocol/pluginsocket.pb.go

.PHONY: lint test dep sync_with_kong
lint: $(proto_def_compiled)
	golangci-lint run --exclude composites

dep: $(proto_def_compiled)
	go get -v
	go mod tidy

test: dep
	go test -v -race ./...

sync_with_kong: clean $(proto_def_compiled)

.PHONY: clean
clean:
	rm -rf $(proto_def)
	rm -rf $(proto_def_compiled)
	

$(proto_def):
	wget https://raw.githubusercontent.com/Kong/kong/master/kong/include/kong/pluginsocket.proto -P $(shell dirname $@)

$(proto_def_compiled): $(proto_def)
	mkdir -p server/kong_plugin_protocol
	protoc -I . $^ --go_out=. --go_opt=paths=source_relative
