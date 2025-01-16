FROM golang:1.23-alpine AS builder

ARG GIT_COMMIT
ARG GIT_TAG_VERSION
ARG BUILD_DATE

WORKDIR /src

ENV CGO_ENABLED="1"
ENV GOARCH="amd64"
ENV GOOS="linux"

COPY . /src/

RUN apk update && apk upgrade && apk add pkgconf git bash build-base sudo && \
        git clone --depth 1 --branch v2.8.0 https://github.com/edenhill/librdkafka.git && \
        cd librdkafka && \
        ./configure --prefix /usr && \
        make && make install

RUN go mod download && \
        go build -tags musl \
        -ldflags="-s -w -X am-kafka-project/pkg/version.Revision=${GIT_COMMIT} -X am-kafka-project/pkg/version.Version=${GIT_TAG_VERSION} -X am-kafka-project/pkg/version.BuildDate=${BUILD_DATE}" \
        -v -o bin/am-kafka-project cmd/am-kafka-project/main.go

FROM alpine:3.21

RUN apk update --no-cache && apk upgrade --no-cache

COPY --from=builder /src/bin/am-kafka-project /usr/bin/am-kafka-project

ENV GIN_MODE="release"
ENV PORT="8080"
ENV AKP_ALERTS_TOPIC="alerts"
ENV AKP_KAFKA_BROKERS="localhost:9091,localhost:9092,localhost:9093"
ENV AKP_KAFKA_CLIENT_ID="gw_alerts"
ENV AKP_SCHEMA_REGISTRY="http://localhost:8081"

EXPOSE 8080

CMD ["/usr/bin/am-kafka-project"]