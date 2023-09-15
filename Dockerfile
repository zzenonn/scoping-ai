FROM golang:1.21 AS builder

RUN mkdir /app
COPY . /app
WORKDIR /app

ENV LOG_LEVEL="DEBUG"

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/main.go

FROM alpine:latest AS prod

COPY --from=builder /app .
CMD ["./app"]