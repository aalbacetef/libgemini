
.PHONY: fmt
fmt:
	goimports -w .

.PHONY: lint 
lint:
	golangci-lint run

.PHONY: test 
test:
	go test -v ./...
