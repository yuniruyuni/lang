FROM alpine:latest@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf

WORKDIR /lang

RUN apk add --update \
    vim \
    bash \
    git \
    make \
    musl-dev \
    go \
    curl \
    llvm \
    clang

ENTRYPOINT ["/bin/bash"]
