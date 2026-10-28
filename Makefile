GO := GOWORK=off go

.PHONY: fmt test vet build lint

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

build:
	$(GO) build ./...

lint:
	golangci-lint run ./...
