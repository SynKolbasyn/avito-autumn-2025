.PHONY: up lint

up:
	docker compose up --build --pull always -d

lint:
	golangci-lint run ./...
