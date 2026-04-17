# build
FROM golang:1.26.2-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# run
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config config/

EXPOSE 8080

CMD ["./server"]
