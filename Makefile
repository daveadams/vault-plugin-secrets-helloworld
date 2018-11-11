PROJECT      = "github.com/daveadams/vault-plugin-secrets-helloworld"
NAME         = $(shell go run version/cmd/main.go name)
VERSION      = $(shell go run version/cmd/main.go version)
COMMIT       = $(shell git rev-parse --short HEAD)
LDFLAGS      = -X ${PROJECT}/version.GitCommit=${COMMIT}
GOFILES      = $(shell find . -name "*.go")

default: vault-plugin-secrets-helloworld

vault-plugin-secrets-helloworld: $(GOFILES)
	@echo "Building ${NAME} ${VERSION} (${COMMIT})"
	go build -ldflags "${LDFLAGS}" ./cmd/vault-plugin-secrets-helloworld

clean:
	rm -f vault-plugin-secrets-helloworld
	rm -rf _workspace

test: vault-plugin-secrets-helloworld
	test/start-server.sh

.PHONY: default clean test
