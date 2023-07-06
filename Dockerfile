FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o /go/bin/crawler

FROM alpine:latest

COPY --from=builder /go/bin/crawler /crawler

RUN chmod +x /crawler

ENV GIN_MODE=release

# rest api
EXPOSE 8989

# debug
EXPOSE 45677

VOLUME [ "/config" ]

ENTRYPOINT [ "/crawler" ]