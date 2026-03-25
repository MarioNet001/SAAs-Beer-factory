# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api/main.go

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
ENV PORT=8080
EXPOSE 8080
CMD ["./main"]
