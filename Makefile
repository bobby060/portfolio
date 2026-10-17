

#  Variables for common things
MODULE_DIR = site/src
BIN_DIR = $(CURDIR)/bin
BINARY = $(BIN_DIR)/site
SRC_DIR = ./...
COMPOSE_PATH = site/compose.yml


.PHONY: generate
generate:
	cd $(MODULE_DIR) && templ generate


.PHONY: build
build: generate
	mkdir -p $(BIN_DIR)
	go -C $(MODULE_DIR) build -o $(BINARY) ./cmd/site

.PHONY: run
run: build
	$(BINARY)

.PHONY: fmt
fmt:
	go -C $(MODULE_DIR) fmt $(SRC_DIR)


.PHONY: test
test: vet
	go -C $(MODULE_DIR) test $(SRC_DIR)

.PHONY: vet
vet:
	go -C $(MODULE_DIR) vet $(SRC_DIR)

.PHONY: build-docker
build-docker: generate
	docker compose -f $(COMPOSE_PATH) build

.PHONY: up
up:
	docker compose  -f $(COMPOSE_PATH) up -d

.PHONY: down
down:
	docker compose -f $(COMPOSE_PATH) down

.PHONY: logs
logs:
	docker compose -f $(COMPOSE_PATH) logs -f

.PHONY: clean
clean:
	rm -r bin/
	find $(MODULE_DIR) -name '*_templ.go' -delete

