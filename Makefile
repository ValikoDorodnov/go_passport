include .env
export

up-project: env-init docker-up print-location

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down -v --remove-orphans

docker-up-db:
	docker-compose up -d --build postgres

env-init:
	rm -f .env && cp -n .env.example .env

print-location:
	@printf "\n   \e[30;42m %s \033[0m\n\n" 'Navigate your browser to ⇒ http://localhost:${REST_ADDR}';

run-test:
	go test ./...