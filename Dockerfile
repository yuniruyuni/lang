FROM alpine:latest

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
