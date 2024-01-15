start:
	DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 DEBUG=true docker-compose build
	docker-compose up

clear: 
	docker-compose down

tests:
	cd notifications-api && \
	go test ./...

	cd notifications-sender && \
	go test ./...