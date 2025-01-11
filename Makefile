#!make

build:
	docker compose up --build

run-migration:
	docker-compose exec app ./aqua-sec-cloud-inventory migrate

seed-db:
	docker-compose exec app ./aqua-sec-cloud-inventory seed

test:
	go test ./... -cover