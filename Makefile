dir = $(shell pwd)

include .env.makefile

stripeContainer = docker run --rm -it --network=go-stripe-eboekhouden-net -e STRIPE_API_KEY=$$STRIPE_API_KEY stripe/stripe-cli:v1.7.9

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
	docker-compose exec api /go-stripe-eboekhouden/go-stripe-eboekhouden migrate up --mysql-conn "mysql://user:password@tcp(db:3306)/db"

.PHONY: status
status:
	docker-compose exec api /go-stripe-eboekhouden/go-stripe-eboekhouden migrate status --mysql-conn "mysql://user:password@tcp(db:3306)/db"

.PHONY: mysql
mysql: 
	docker-compose exec db mysql -u root --password=password --database=db

.PHONY: listen
listen:
	$(stripeContainer) listen -f api:8080/hooks