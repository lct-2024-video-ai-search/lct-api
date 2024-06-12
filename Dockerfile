# Build stage
FROM golang:1.22-alpine3.18 AS builder
ENV GOGACHE="/go/cache"

WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
#FROM alpine:3.18
#
## Copy config file
#COPY app.env .
#
## Copy migrations
#COPY db/migrations ./db/migrations
#
## Copy binary and start script
#COPY start.sh .
#COPY --from=builder /app/main .
#RUN chmod +x start.sh main
#
#EXPOSE 8080
#ENTRYPOINT ["./start.sh"]

# Acts as default argument for ENTRYPOINT above
CMD ["./main"]

