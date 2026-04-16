#!/bin/bash

DATABASE_URL="postgres://myuser:mypassword@host.docker.internal:5432/mydatabase?sslmode=disable"
MIGRATIONS_DIR="./migrations"

case "$1" in
  up)
    docker run --rm -v $(pwd)/$MIGRATIONS_DIR:/migrations \
      migrate/migrate -path=/migrations -database $DATABASE_URL up
    ;;
  down)
    docker run --rm -v $(pwd)/$MIGRATIONS_DIR:/migrations \
      migrate/migrate -path=/migrations -database $DATABASE_URL down ${2:-1}
    ;;
  create)
    docker run --rm -v $(pwd)/$MIGRATIONS_DIR:/migrations \
      migrate/migrate create -ext sql -dir /migrations -seq "$2"
    ;;
  version)
    docker run --rm -v $(pwd)/$MIGRATIONS_DIR:/migrations \
      migrate/migrate -path=/migrations -database $DATABASE_URL version
    ;;
  *)
    echo "Usage: $0 {up|down|create|version} [name|count]"
    exit 1
esac
