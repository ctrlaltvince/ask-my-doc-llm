# backend/Dockerfile
FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Create non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser
COPY --from=builder /app/app .

RUN chown appuser:appgroup /home/appuser/app

USER appuser

EXPOSE 8081
CMD ["./app"]
