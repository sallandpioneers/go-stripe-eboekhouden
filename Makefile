dir = $(shell pwd)

include .env.makefile

stripeContainer = docker run --rm -it --network=go-stripe-eboekhouden-net -e STRIPE_API_KEY=$$STRIPE_API_KEY stripe/stripe-cli:v1.8.0

.PHONY: default
default:
	-docker network create go-stripe-eboekhouden-net

.PHONY: build
build:
	docker-compose build

.PHONY: up
up: default
	docker-compose up -d --remove-orphans

.PHONY: logs
logs:
	docker-compose logs -f

.PHONY: down
down:
	docker-compose down --remove-orphans

.PHONY: test
test: 
	docker exec -it go-stripe-eboekhouden_api_1 go test ./... -v -cover

.PHONY: bench
bench: 
	docker exec -it go-stripe-eboekhouden_api_1 go test ./... -bench=.

.PHONY: migrate
migrate:
	docker-compose exec go_stripe_boekhouden_api /go-stripe-eboekhouden/go-stripe-eboekhouden migrate up --mysql-conn "mysql://user:password@tcp(go_stripe_boekhouden_db:3306)/db"

.PHONY: status
status:
	docker-compose exec go_stripe_boekhouden_api /go-stripe-eboekhouden/go-stripe-eboekhouden migrate status --mysql-conn "mysql://user:password@tcp(go_stripe_boekhouden_db:3306)/db"

.PHONY: mysql
mysql: 
	docker-compose exec go_stripe_boekhouden_db mysql -u root --password=password --database=db

.PHONY: sqlc
sqlc:
	rm -f internal/storage/mysql/queries/generated/*.sql.go
	docker run --rm -v $(dir)/internal/storage/mysql:/src -w /src/queries kjconroy/sqlc generate

.PHONY: listen
listen:
	$(stripeContainer) listen -f go_stripe_boekhouden_api:8080/hooks