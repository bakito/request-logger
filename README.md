# Request Logger

A small go webserver logging all requests

## Run

```bash
docker run -p 8080:8080 bakito/request-logger
```

## Output

The output is simply generated with https://golang.org/pkg/net/http/httputil/#DumpRequest

The request number is returned as response header 'Request-No'.

```bash

2018/11/08 23:13:59 Request-No: 5
POST / HTTP/1.1
Host: localhost:8080
Accept: */*
Content-Length: 1
Content-Type: application/x-www-form-urlencoded
User-Agent: curl/7.61.1

<body>
```