
.PHONY: fmt
fmt:
	goimports -w .

.PHONY: lint 
lint:
	golangci-lint run
