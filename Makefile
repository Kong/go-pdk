
.PHONY: lint
lint:
	golint

.PHONY: dep
dep:
	go mod tidy

.PHONY: test
test:
	go test -v -race ./...
	
.PHONY: verify-protoc-gen
verify-protoc-gen:
	./scripts/verify-protoc-gen.sh

.PHONY: update-protoc-gen
update-protoc-gen:
	./scripts/update-protoc-gen.sh
