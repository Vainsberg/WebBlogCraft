run:
	go run cmd/server/main.go

start:
	docker-compose up -d
	
stop:
	docker-compose stop