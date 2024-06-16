# Build stage
FROM golang:1.22-alpine3.18 AS builder
ENV GOGACHE="/go/cache"

WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.18

# Copy config file
COPY app.env .

# Copy binary and start script
COPY --from=builder /app/main .
EXPOSE 8080

ENTRYPOINT ["./main"]

