FROM golang:1.24.3-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o main ./cmd/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main /app/main
EXPOSE 8080
CMD ["/app/main"]