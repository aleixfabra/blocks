.PHONY: up down restart logs blocks mempool

up:
	docker compose -f blocks/docker-compose.yaml up -d
	docker compose -f mempool/docker-compose.yaml up

down:
	docker compose -f mempool/docker-compose.yaml down
	docker compose -f blocks/docker-compose.yaml down

restart:
	$(MAKE) down
	$(MAKE) up