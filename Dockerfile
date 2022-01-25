FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /go/bin/crawler

FROM alpine:latest

COPY --from=builder /go/bin/crawler /crawler
RUN chmod +x /blive

ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT [ "/crawler" ]