FROM --platform=${BUILDPLATFORM} curlimages/curl:7.85.0 as build-stage0

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

ENV KP_VERSION 0.7.0

RUN curl -L -o /tmp/kp https://github.com/vmware-tanzu/kpack-cli/releases/download/v${KP_VERSION}/kp-linux-${KP_VERSION} && \
	chmod a+x /tmp/kp

RUN curl -L -o /tmp/kubectl https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
	chmod a+x /tmp/kubectl

FROM --platform=${BUILDPLATFORM} ubuntu:20.04
LABEL maintainer="Gareth Evans <gareth@bryncynfelin.co.uk>"

COPY --from=build-stage0 /tmp/kp /usr/bin/kp
COPY --from=build-stage0 /tmp/kubectl /usr/bin/kubectl

COPY github-actions-entrypoint.sh /usr/bin/github-actions-entrypoint.sh

RUN kp --help

ENTRYPOINT [ "/usr/bin/github-actions-entrypoint.sh" ]
