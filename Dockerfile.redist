# Build stage
FROM golang:1.10 as builder

WORKDIR /usr/bin/
RUN curl -sLSf https://raw.githubusercontent.com/alexellis/license-check/master/get.sh | sh

WORKDIR /go/src/github.com/openfaas-incubator/ofc-bootstrap
COPY . .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }


RUN /usr/bin/license-check -path ./ --verbose=false "Alex Ellis" "OpenFaaS Author(s)" \
 && go test $(go list ./... | grep -v /vendor/ | grep -v /template/|grep -v /build/) -cover \
 && VERSION=$(git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///') \
 && GIT_COMMIT=$(git rev-list -1 HEAD) \
 && CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w \
        -X github.com/openfaas-incubator/ofc-bootstrap/version.GitCommit=${GIT_COMMIT} \
        -X github.com/openfaas-incubator/ofc-bootstrap/version.Version=${VERSION}" \
        -a -installsuffix cgo -o ofc-bootstrap \
 && CGO_ENABLED=0 GOOS=darwin go build --ldflags "-s -w \
        -X github.com/openfaas-incubator/ofc-bootstrap/version.GitCommit=${GIT_COMMIT} \
        -X github.com/openfaas-incubator/ofc-bootstrap/version.Version=${VERSION}" \
        -a -installsuffix cgo -o ofc-bootstrap-darwin \
 && CGO_ENABLED=0 GOOS=windows go build --ldflags "-s -w \
        -X github.com/openfaas-incubator/ofc-bootstrap/version.GitCommit=${GIT_COMMIT} \
        -X github.com/openfaas-incubator/ofc-bootstrap/version.Version=${VERSION}" \
        -a -installsuffix cgo -o ofc-bootstrap.exe

# Release stage
FROM alpine:3.8

RUN apk --no-cache add ca-certificates git

WORKDIR /root/

COPY --from=builder /go/src/github.com/openfaas-incubator/ofc-bootstrap/ofc-bootstrap                .
COPY --from=builder /go/src/github.com/openfaas-incubator/ofc-bootstrap/ofc-bootstrap-darwin         .
COPY --from=builder /go/src/github.com/openfaas-incubator/ofc-bootstrap/ofc-bootstrap.exe            .

ENV PATH=$PATH:/root/

CMD ["ofc-bootstrap"]
