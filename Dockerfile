FROM golang:1.13-alpine as go_builder
WORKDIR /go/src/blog-api
COPY . .
RUN apk add --no-cache git
RUN apk add curl && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -o blog-api

FROM scratch
COPY --from=go_builder /go/src/blog-api/blog-api .
ENV BASE_URL=0.0.0.0
EXPOSE 8080
CMD ["./blog-api"]
