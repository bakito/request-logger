FROM quay.io/bitnami/golang:1.15 as builder

WORKDIR /build

RUN apt-get update && apt-get install -y upx
COPY . .

ENV GOPROXY=https://goproxy.io \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN go build -a -installsuffix cgo -ldflags="-w -s" -o request-logger && \
    upx --ultra-brute -q request-logger

# application image

FROM scratch

LABEL maintainer="bakito <github@bakito.ch>"
EXPOSE 8080
USER 1001
ENTRYPOINT ["/go/bin/request-logger"]

COPY --from=builder /build/request-logger /go/bin/request-logger
