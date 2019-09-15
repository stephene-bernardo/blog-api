FROM golang:1.12-alpine

WORKDIR /go/src/blog-api
COPY . .
RUN apk add curl && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && dep ensure
CMD CGO_ENABLED=0 go test ./...