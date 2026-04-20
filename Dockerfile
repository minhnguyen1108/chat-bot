# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot .

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bot .
RUN apk add --no-cache ca-certificates
ENV PORT=8080
EXPOSE 8080
CMD ["./bot"]
