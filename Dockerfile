FROM golang:1.24-alpine AS builder

WORKDIR /build

RUN apk update && apk add upx
COPY . .

ENV GOPROXY=https://goproxy.io \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux
RUN go build -a -installsuffix cgo -ldflags="-w -s" -o request-logger && \
    upx -q request-logger

# application image

FROM scratch

LABEL maintainer="bakito <github@bakito.ch>"
EXPOSE 8080
USER 1001
ENTRYPOINT ["/go/bin/request-logger"]

COPY --from=builder /build/request-logger /go/bin/request-logger
