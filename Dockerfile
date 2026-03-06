# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /server .

# Run stage
FROM alpine:3.21

RUN apk --no-cache add ca-certificates
WORKDIR /

COPY --from=builder /server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
