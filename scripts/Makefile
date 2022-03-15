DEV_PROJECT := pow-f
DEV_COMPOSE_FILE := scripts/docker-compose.yml

# -----------------------------------------------------------------------------
.PHONY: app
app:
	${INFO} "Removing old instance app (if any)"
	@docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) stop go-front
	@docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) rm -f go-front
	${INFO} "Running backend application"
	@docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) up -d go-front
