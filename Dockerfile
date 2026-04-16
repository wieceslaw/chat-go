# TODO: fix
FROM golang:1.21-alpine

WORKDIR /app
COPY . .

# Build with config path
RUN go build -o server ./cmd/server

# Use environment variable for config
ENV APP_ENV=production
ENV DB_PASSWORD=from_secret

EXPOSE 8080
CMD ["./server", "-config", "/etc/myapp/config.prod.yaml"]
