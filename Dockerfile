FROM --platform=${BUILDPLATFORM} golang:1.19 AS build-stage0

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

ARG VERSION=unknown

WORKDIR /root/
COPY . ./

RUN go build -trimpath -ldflags "-X github.com/garethjevans/build-action/pkg/version.Version=$VERSION" -o builder main.go

FROM --platform=${BUILDPLATFORM} ubuntu:20.04
LABEL maintainer="Gareth Evans <gareth@bryncynfelin.co.uk>"

COPY --from=build-stage0 /root/builder /usr/bin/builder
COPY github-actions-entrypoint.sh /usr/bin/github-actions-entrypoint.sh

ENTRYPOINT [ "/usr/bin/github-actions-entrypoint.sh" ]
