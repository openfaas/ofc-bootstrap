
GO_FILES?=$$(find . -name '*.go' |grep -v vendor)
TAG?=latest

.PHONY: build

install-ci:
	./hack/install-ci.sh
ci:
	./hack/ci.sh

.PHONY: build
build:
	./build.sh

.PHONY: dist
dist:
	./build_redist.sh
