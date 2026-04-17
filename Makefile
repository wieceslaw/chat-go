run-local:
	# todo: run db + migrations
	go run cmd/server/main.go

docker-run:
	docker-compose up -d

docker-clean:
	docker-compose down -v

docker-build-debug:
	docker build -t chatapp:latest . --no-cache --progress=plain
