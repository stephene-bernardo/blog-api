FROM golang:1.13-alpine as go_builder
WORKDIR /go/src/blog-api
COPY . .
RUN apk add curl && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -o blog-api

FROM postgres:10-alpine
ENV BASE_URL=0.0.0.0 POSTGRES_HOST=localhost POSTGRES_PASSWORD=abc123 POSTGRES_DB=blog
COPY ./entrypoint.sh .
COPY --from=go_builder /go/src/blog-api/blog-api .
EXPOSE 8080 5422
CMD ["./entrypoint.sh"]
