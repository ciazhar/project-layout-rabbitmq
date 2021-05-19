compose:
	docker compose up -d

run-subscriber:
	go run cmd/subscriber/main.go

run-producer:
	go run cmd/producer/main.go