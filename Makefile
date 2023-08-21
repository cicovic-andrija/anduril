# Makefile for Go module github.com/cicovic-andrija/anduril
#

OUTPUT_DIR = out
DATA_DIR = $(OUTPUT_DIR)/data
SERVER_BIN = $(OUTPUT_DIR)/anduril-server
MKCONF = $(OUTPUT_DIR)/mkconf
CONFIG_TEMPLATE = configuration/anduril-config.json
CONFIG_FILE = $(DATA_DIR)/encrypted-config.txt

.PHONY: build
build: css | $(OUTPUT_DIR)
	go build -o $(SERVER_BIN) main.go

.PHONY: css
css:
	sass assets/stylesheets/main.scss:assets/stylesheets/styles.css

.PHONY: config
config: tools
	$(MKCONF) \
		--template $(CONFIG_TEMPLATE) \
		--to $(CONFIG_FILE) \
		--profile prod \
		--decrypt | jq

.PHONY: tools | $(OUTPUT_DIR)
tools: | $(OUTPUT_DIR)
	go build -o $(MKCONF) ./tools/mkconf.go

.PHONY: all
all: $(OUTPUT_DIR) build tools config

.PHONY: devenv
devenv: build tools
	mkdir -p $(DATA_DIR)/assets
	rsync -rv assets/templates $(DATA_DIR)/
	rsync -v assets/scripts/*.js $(DATA_DIR)/assets/
	rsync -v assets/stylesheets/*.css $(DATA_DIR)/assets/
	rsync -rv assets/icons $(DATA_DIR)/assets/
	$(MKCONF) \
		--template $(CONFIG_TEMPLATE) \
		--to $(CONFIG_FILE) \
		--profile dev \
		--decrypt | jq
	openssl req \
		-x509 \
		-newkey rsa:4096 \
		-sha256 \
		-days 365 \
		-nodes \
		-out $(OUTPUT_DIR)/tlspublic.crt \
		-keyout $(OUTPUT_DIR)/tlsprivate.key \
		-subj "/CN=localhost/C=/ST=/L=/O=/OU=" \
		>/dev/null 2>/dev/null

$(OUTPUT_DIR):
	mkdir -p $(DATA_DIR)

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)/logs
	rm -rf $(OUTPUT_DIR)/work
	rm -f $(OUTPUT_DIR)/tls*
	rm -rf $(DATA_DIR)
	rm -f $(SERVER_BIN)
	rm -f $(MKCONF)

.PHONY: tidy
tidy:
	go mod tidy
