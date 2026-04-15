# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
ARG REGISTRY="docker.io/library"
ARG BASE_IMAGE=golang
ARG BASE_TAG=1.26-trixie@sha256:c0074c718b473f3827043f86532c4c0ff537e3fe7a81b8219b0d1ccfcc2c9a09

FROM $REGISTRY/$BASE_IMAGE:$BASE_TAG AS builder
ENV DEBIAN_FRONTEND=noninteractive
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on GOPATH=/src/tmp/go

ARG XDG_CONFIG_HOME

RUN apt-get update && \
    apt-get upgrade -y
RUN git config --global url."git@github.com:AustralianCyberSecurityCentre/".insteadOf "https://github.com/AustralianCyberSecurityCentre/"

# copy somewhere outside GOPATH/src as using go modules
COPY . /src

# full static builds with no ld deps, so we can copy it to scratch
RUN cd /src && go build -v -a -tags netgo -ldflags '-w -extldflags "-static"' -o /go/bin/azul-entropy *.go

# now copy artifacts to a lightweight image
FROM $REGISTRY/alpine:latest@sha256:c69a6ff7c24d1ffa913798501d0e7104e0e9764e28eb44a930939f91ef829e64
COPY --from=builder /go/bin /bin
ENTRYPOINT ["/bin/azul-entropy"]
