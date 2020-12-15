Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X github.com/openfaas/ofc-bootstrap/cmd.Version=$(Version) -X github.com/openfaas/ofc-bootstrap/cmd.GitCommit=$(GitCommit)"
SOURCE_DIRS = cmd pkg main.go
export GO111MODULE=on

.PHONY: all
all: gofmt test dist hash

.PHONY: ci
ci: all install-ci ci

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o ofc-bootstrap

.PHONY: gofmt
gofmt:
	@test -z $(shell gofmt -l -s $(SOURCE_DIRS) ./ | tee /dev/stderr) || (echo "[WARN] Fix formatting issues with 'make fmt'" && exit 1)

.PHONY: test
test:
	CGO_ENABLED=0 go test $(shell go list ./... | grep -v /vendor/|xargs echo) -cover

.PHONY: dist
dist:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/ofc-bootstrap
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/ofc-bootstrap-darwin
	CGO_ENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/ofc-bootstrap.exe

.PHONY: hash
hash:
	rm -rf bin/*.sha256 && ./hack/hashgen.sh

.PHONY: install-ci
install-ci:
	./hack/install-ci.sh

.PHONY: ci
ci:
	./hack/integration-test.sh
