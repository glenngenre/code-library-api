FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate swagger docs before building
RUN swag init -g main.go --output docs

RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ---

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 3000

CMD ["./server"]
