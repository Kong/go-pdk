
.PHONY: lint
lint:
	golangci-lint run --exclude composites

.PHONY: dep
dep:
	go get -v
	go mod tidy

.PHONY: test
test:
	go test -v -race ./...