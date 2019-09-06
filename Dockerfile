FROM golang:1.13 as builder

WORKDIR /go/src/github.com/bakito/request-logger

COPY ./* /go/src/github.com/bakito/request-logger/

RUN apt-get update && apt-get install -y xz-utils && \
  curl -SL --fail --silent --show-error https://github.com/upx/upx/releases/download/v3.95/upx-3.95-amd64_linux.tar.xz | tar --wildcards -xJ --strip-components 1 */upx

ENV GOPROXY=https://goproxy.io \ 
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -a -installsuffix cgo -ldflags="-w -s" -o request-logger && \
  ./upx --ultra-brute -q request-logger

# application image

FROM scratch

LABEL maintainer="bakito <github@bakito.ch>"
EXPOSE 8080
USER 1001
ENTRYPOINT ["/go/bin/request-logger"]

COPY --from=builder /go/src/github.com/bakito/request-logger/request-logger /go/bin/request-logger
