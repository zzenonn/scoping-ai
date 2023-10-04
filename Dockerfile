# Stage 1: Builder
FROM golang:1.21 AS builder

RUN mkdir /app
COPY . /app
WORKDIR /app

ENV LOG_LEVEL="INFO"

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/main.go

# Stage 2: Production image
FROM alpine:latest AS prod

# Copying from builder stage
COPY --from=builder /app .

# Using the project ID as a flag for running the app
CMD ["./app", "--project-id", "genai-tna"]
