
.PHONY: lint
lint:
	golint

.PHONY: dep
dep:
	go mod tidy

.PHONY: test
test:
	go test -v -race ./...
	
.PHONY: coverage
coverage:
	go test -race -v -count=1 -coverprofile=coverage.out ./...	
