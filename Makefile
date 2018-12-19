.PHONY: ci

install-ci:
	./hack/install-ci.sh
ci:
	./hack/ci.sh
