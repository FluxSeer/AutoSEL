GO := GOWORK=off go

.PHONY: fmt test vet build lint selinux-qualification

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

selinux-qualification:
	./scripts/selinux/qualification.sh
