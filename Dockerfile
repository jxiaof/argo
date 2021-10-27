FROM golang:1.14.2-alpine3.11 as builder

WORKDIR /go/src

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update \
    && rm -rf /var/cache/apk/*


ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
COPY src/go.mod go.mod
RUN go mod download
COPY src/ .
RUN go build -o server .


FROM alpine:latest as prod

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update && apk upgrade \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates 2>/dev/null || true \
    && rm -rf /var/cache/apk/


ENV PROJ_NAME=iris-web
ENV PROJ_VERSION=v1
WORKDIR /go/src/$PROJ_NAME

COPY --from=0 /go/src/server .

EXPOSE 8090

LABEL name="iris-web"
LABEL version="v1"

CMD "/go/src/$PROJ_NAME/server"