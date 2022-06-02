FROM golang:1.16.15-alpine3.15 AS builder
LABEL AUTHOR Seungkyu Ahn (seungkyua@gmail.com)

ENV KUSTOMIZE_VER v4.2.0
ENV HOME /root
ENV GOPATH $HOME/golang
# on means using vendor directory
ENV GO111MODULE on
# CGO_ENABLED=1 default value means using link
ENV CGO_ENABLED 1


RUN apk update && apk add git curl tar bash build-base


WORKDIR $HOME
COPY . $HOME/kustomize-helm-transformer


# install kustomize from source
WORKDIR $HOME
RUN git clone https://github.com/kubernetes-sigs/kustomize.git
WORKDIR $HOME/kustomize
RUN git checkout -b tags_${KUSTOMIZE_VER} tags/kustomize/${KUSTOMIZE_VER}
WORKDIR $HOME/kustomize/kustomize
RUN unset GOPATH && unset GO111MODULES && go install .
RUN cp $HOME/go/bin/kustomize /usr/local/bin/


# plugin copy
RUN rm -rf $HOME/kustomize/plugin/*
RUN cp -r $HOME/kustomize-helm-transformer/plugin/openinfradev.github.com $HOME/kustomize/plugin/


WORKDIR $HOME/kustomize/plugin/openinfradev.github.com/v1/helmvaluestransformer
RUN rm -rf vendor && rm -f go.mod go.sum
RUN mv kustomize-${KUSTOMIZE_VER}-go.mod go.mod
# update go.mod and go.sum
RUN unset GOPATH && unset GO111MODULES && go mod tidy


WORKDIR $HOME/kustomize
RUN unset GOPATH && unset GO111MODULES && ./hack/buildExternalGoPlugins.sh ./plugin





FROM alpine:edge
LABEL AUTHOR Seungkyu Ahn (seungkyua@gmail.com)

RUN apk add --no-cache bash git

USER root

RUN mkdir -p /root/.config/kustomize/plugin/openinfradev.github.com/v1/helmvaluestransformer
COPY --from=builder /root/kustomize/plugin/openinfradev.github.com/v1/helmvaluestransformer/HelmValuesTransformer.so /root/.config/kustomize/plugin/openinfradev.github.com/v1/helmvaluestransformer/
COPY --from=builder /usr/local/bin/kustomize /usr/local/bin/kustomize
WORKDIR /root

CMD ["kustomize"]