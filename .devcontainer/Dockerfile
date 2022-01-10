# [Choice] Ubuntu version: bionic, focal
ARG VARIANT="bionic"
FROM mcr.microsoft.com/vscode/devcontainers/base:0-${VARIANT}

# install common packages
RUN apt-get update && apt install -y \
    software-properties-common \
    gcc

# install go
ARG GO_VERSION=1.17.6
RUN wget https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz && \
    sha256sum go$GO_VERSION.linux-amd64.tar.gz && \
    tar -xvf go$GO_VERSION.linux-amd64.tar.gz && \
    mv go /usr/local && \
    rm go$GO_VERSION.linux-amd64.tar.gz

# set up env vars
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH="$GOPATH/bin:$GOROOT/bin:$PATH"

# install go deps
RUN go get -x \
        github.com/spf13/cobra/cobra \
        github.com/spf13/cobra \
        github.com/spf13/viper \
        github.com/ramya-rao-a/go-outline \
        github.com/stamblerre/gocode \
        github.com/uudashr/gopkgs/v2/cmd/gopkgs \
        golang.org/x/tools/cmd/goimports \
        github.com/cweill/gotests/gotests \
        golang.org/x/tools/cmd/guru

# gopls is a special flower and does not like to be installed with the others
RUN go get -x golang.org/x/tools/gopls

# install debugger
RUN GO111MODULE=on GOBIN=/tmp/ go get github.com/go-delve/delve/cmd/dlv@master && \
    mv /tmp/dlv $GOPATH/bin/dlv-dap

