PORT=15521
VOLUME_NAME=mailpocket-data

.PHONY: run-batched run-sqlite stop reset

test-vol:
	@if ! docker volume ls -q | grep -q "^$(VOLUME_NAME)$$"; then \
		echo "Creating volume $(VOLUME_NAME)..."; \
		docker volume create $(VOLUME_NAME); \
	fi

run-batched: test-vol stop
	docker build -t batched-server -f Dockerfile.batched .
	docker run -p $(PORT):$(PORT) --name batched-server -v $(VOLUME_NAME):/app/data batched-server

run-sqlite: test-vol stop
	docker build -t sqlite-server -f Dockerfile.sqldb .
	docker run -d -p $(PORT):$(PORT) --name sqlite-server -v $(VOLUME_NAME):/app/data sqlite-server

stop:
	docker stop batched-server || true
	docker stop sqlite-server || true
	docker rm batched-server || true
	docker rm sqlite-server || true

reset: stop
	docker volume rm $(VOLUME_NAME) || true