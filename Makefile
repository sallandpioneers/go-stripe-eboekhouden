dir = $(shell pwd)

.PHONY: build
build:
	docker-compose build

.PHONY: up
up:
	docker-compose up -d --remove-orphans

.PHONY: logs
logs:
	docker-compose logs -f

.PHONY: down
down:
	docker-compose down --remove-orphans

.PHONY: test
test: 
	docker exec -it stripe-eboekhouden-go_api_1 go test ./... -v -cover

.PHONY: bench
bench: 
	docker exec -it stripe-eboekhouden-go_api_1 go test ./... -bench=.

