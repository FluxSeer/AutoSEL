GO := GOWORK=off go

.PHONY: fmt test vet build lint selinux-qualification docker-readonly-qualification

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

docker-readonly-qualification:
	./scripts/docker/qualification.sh
