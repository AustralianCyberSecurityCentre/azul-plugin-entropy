# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
ARG REGISTRY="docker.io/library"
ARG BASE_IMAGE=golang
ARG BASE_TAG=1.27-trixie@sha256:2f84bc93ecfb2689f782b153fdcd368b5a7ab96c1386c65cdaccf35e726d6a44

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
FROM $REGISTRY/alpine:latest@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6
COPY --from=builder /go/bin /bin
ENTRYPOINT ["/bin/azul-entropy"]
