# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
ARG REGISTRY="docker.io/library"
ARG BASE_IMAGE=golang
ARG BASE_TAG=1.27-trixie@sha256:af1bfc2614ebb1dbd31af514b78740fd678850e76498762b425ecee4c5c6ea80

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
FROM $REGISTRY/alpine:latest@sha256:3f4d6b32e26703c6082839234f43769d5046b3ac7c910033911d1e71d8cd630a
COPY --from=builder /go/bin /bin
ENTRYPOINT ["/bin/azul-entropy"]
