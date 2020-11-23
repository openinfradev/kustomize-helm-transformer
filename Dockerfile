FROM alpine:latest

MAINTAINER SKTelecom 5GX Cloud Labs

ENV HOME /root
ENV GO111MODULE on
ENV PATH /usr/local/go/bin:$PATH
ENV GOROOT /usr/local/go
ENV GOPATH $HOME/golang

RUN mkdir -p $HOME/golang
RUN mkdir -p $HOME/.config/kustomize/plugin/openinfradev.github.com/v1/helmvaluestransformer

RUN apk update && apk add --no-cache curl git jq openssh libc6-compat build-base

WORKDIR $HOME
RUN git clone https://github.com/openinfradev/kustomize-helm-transformer.git
RUN cat kustomize-helm-transformer/README.md | grep -m 1 "* kustomize" | sed -nre 's/^[^0-9]*(([0-9]+\.)*[0-9]+).*/\1/p' > .kustomize_version
RUN cat kustomize-helm-transformer/README.md | grep -m 1 "* go" | sed -nre 's/^[^0-9]*(([0-9]+\.)*[0-9]+).*/\1/p' > .golang_version

WORKDIR /usr/local
RUN curl -fL https://dl.google.com/go/go$(cat $HOME/.golang_version).linux-amd64.tar.gz | tar xz

RUN go get sigs.k8s.io/kustomize/kustomize/v3@v$(cat $HOME/.kustomize_version)
RUN mv $GOPATH/bin/kustomize /usr/local/bin/

WORKDIR $HOME/kustomize-helm-transformer/plugin/openinfradev.github.com/v1/helmvaluestransformer/
RUN go test
RUN mv HelmValuesTransformer.so $HOME/.config/kustomize/plugin/openinfradev.github.com/v1/helmvaluestransformer/

WORKDIR /
