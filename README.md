[![Go](https://github.com/bakito/request-logger/actions/workflows/go.yml/badge.svg)](https://github.com/bakito/request-logger/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bakito/request-logger)](https://goreportcard.com/report/github.com/bakito/request-logger)

# Request Logger

A small go webserver logging all requests

## Run

```bash
docker run -p 8080:8080 ghcr.io/bakito/request-logger
```

## Paths

### Echo

To get the request also visualized in the response call any url with sub-path **/echo**. This retuns the request dump values also in the response.

### Custom Response Code

To get a custom response code (2xx, 4xx or 5xx) call any url with sub path **/code/{code}**.

- code: the response code to be returned

### Random Response Code

To get a response code on a random basis call any url with sub path **/random/code/{code}/{percentage}**.

- code: the response code to be returned at random requests
- percentage: the percentage the response code should be returned, where 1 is always and 0 is never

### Random Sleep

To get a response wit a random delay call any url with sub path **/random/sleep/{sleep}**.

- sleep: the may amount of milliseconds to sleep

### Replay Mode

Predefined responses (content-type and body) can be returned on the path **/replay**.  
To train a response, send a request to the same url with the content-type and body to be returned later. Set the request header **Train-Replay** to **true** for the training request.  
The all following request to this url will get the same content-type and body as it was trained before.


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

## Setup on Openshift

```bash
oc process -f openshift\openshift-template.yaml -p NAME=request-logger | oc apply -f -

```
