# Build stage
FROM golang:1.13 as builder
ENV GO111MODULE=off

ENV CGO_ENABLED=0

WORKDIR /usr/bin/
RUN curl -sLSf https://raw.githubusercontent.com/alexellis/license-check/master/get.sh | sh

WORKDIR /go/src/github.com/openfaas-incubator/ofc-bootstrap
COPY . .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }

RUN /usr/bin/license-check -path ./ --verbose=false "Alex Ellis" "OpenFaaS Author(s)"
RUN go test $(go list ./... | grep -v /vendor/ | grep -v /template/|grep -v /build/) -cover \
    && VERSION=$(git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///') \
    && GIT_COMMIT=$(git rev-list -1 HEAD) \
    && CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w \
    -X github.com/openfaas-incubator/ofc-bootstrap/version.GitCommit=${GIT_COMMIT} \
    -X github.com/openfaas-incubator/ofc-bootstrap/version.Version=${VERSION}" \
    -a -installsuffix cgo -o ofc-bootstrap

# Release stage
FROM alpine:3.11

RUN apk --no-cache add ca-certificates git

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk add --no-cache ca-certificates

WORKDIR /home/app

COPY --from=builder /go/src/github.com/openfaas-incubator/ofc-bootstrap/ofc-bootstrap               /usr/bin/
RUN chown -R app:app ./

USER app

ENV PATH=$PATH:/usr/bin/

CMD ["ofc-bootstrap"]
