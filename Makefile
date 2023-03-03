.PHONY: build
build:
	@sqlc generate
	@echo "Building..."
	@docker build -t authorservice:local .

.PHONY: clean
clean:
	@rm -f ./pkg/database/db.go ./pkg/database/models.go ./pkg/database/queries.sql.go
	@rm authorservice
	@echo "Cleaning..."
	@docker rmi authorservice:local


.PHONY: stop
stop:
	@echo "Stopping..."
	@docker compose -f stack.yml down -v

.PHONY: run
run:
	@echo "Running..."
	@docker compose -f stack.yml down -v
	@docker compose -f stack.yml up

.PHONY: test
test:
	@echo "Testing..."
	@docker compose -f test/docker-compose.yml down -v
	@docker compose -f test/docker-compose.yml up --build --abort-on-container-exit --remove-orphans --force-recreate
	@docker compose -f test/docker-compose.yml down -v