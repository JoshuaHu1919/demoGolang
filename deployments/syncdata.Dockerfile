FROM golang:1.17.2 AS builder

WORKDIR /syncdata
COPY . .

ENV GO111MODULE=on

WORKDIR /syncdata/cmd/worker
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -o syncdataService

FROM alpine:latest
RUN apk update && \
    apk upgrade && \
    apk add --no-cache tzdata  && \
    apk add --no-cache ca-certificates && \
    apk add --no-cache curl && \
    rm -rf /var/cache/apk/*
ARG env
WORKDIR /syncdata
COPY --from=builder /syncdata/cmd/worker/syncdataService /syncdata/
COPY --from=builder /syncdata/deployments/config.${env:-dev}.yml /syncdata/configs/config.yml

CMD ./syncdataService