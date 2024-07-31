.PHONY: build up down

build:
	docker build .

up:
	docker-compose up

down:
	docker-compose down

