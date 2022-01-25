FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o /go/bin/crawler

FROM alpine:latest

COPY --from=builder /go/bin/crawler /crawler
RUN chmod +x /crawler

ENV GIN_MODE=release

EXPOSE 8989

VOLUME [ "/config" ]

ENTRYPOINT [ "/crawler" ]