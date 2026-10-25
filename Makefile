# Variables for common things
MODULE_DIR = site/src
BIN_DIR = $(CURDIR)/bin
BINARY = $(BIN_DIR)/site
SRC_DIR = ./...
COMPOSE_PATH = site/compose.yml

# Tailwind standalone CLI (no Node). Bump the version and checksum together, here and in site/Dockerfile.
TAILWIND_VERSION = v4.3.3
TAILWIND_SHA256 = dc61b3ac6b8c9ca874c0cc4c57b2409791a64c5540404ca5f5367360babc313a
TAILWIND = $(BIN_DIR)/tailwindcss-$(TAILWIND_VERSION)
CSS_IN = $(MODULE_DIR)/web/static/css/input.css
CSS_OUT = $(MODULE_DIR)/web/static/css/site.css


.PHONY: generate
generate:
	go -C $(MODULE_DIR) tool templ generate

# Download Tailwind once per version and verify its checksum before using it
$(TAILWIND):
	mkdir -p $(BIN_DIR)
	curl -fsSL -o $@.tmp https://github.com/tailwindlabs/tailwindcss/releases/download/$(TAILWIND_VERSION)/tailwindcss-linux-x64
	echo "$(TAILWIND_SHA256)  $@.tmp" | sha256sum -c -
	chmod +x $@.tmp && mv $@.tmp $@

.PHONY: css
css: $(TAILWIND)
	$(TAILWIND) -i $(CSS_IN) -o $(CSS_OUT) --minify

.PHONY: build
build: generate css
	mkdir -p $(BIN_DIR)
	go -C $(MODULE_DIR) build -o $(BINARY) ./cmd/site

.PHONY: run
run: build
	$(BINARY)

# Live reload: templ proxy on http://localhost:7331, Tailwind rebuilding CSS on change
.PHONY: dev
dev:
	$(MAKE) -j2 dev-templ dev-css

.PHONY: dev-templ
dev-templ:
	cd $(MODULE_DIR) && CONTENT_DIR=content go tool templ generate --watch \
		--watch-pattern='(.+\.go$$)|(.+\.templ$$)|(.+_templ\.txt$$)|(.+site\.css$$)' \
		--proxy=http://localhost:8080 --cmd="go run ./cmd/site"

.PHONY: dev-css
dev-css: $(TAILWIND)
	$(TAILWIND) -i $(CSS_IN) -o $(CSS_OUT) --watch

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
build-docker:
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
	rm -rf $(BIN_DIR) $(CSS_OUT)
	find $(MODULE_DIR) -name '*_templ.go' -delete
