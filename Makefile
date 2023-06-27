# Makefile for Go module github.com/cicovic-andrija/anduril
#

OUTPUT_DIR = out
SERVER_BIN = anduril-server
MKCONF_PATH = $(OUTPUT_DIR)/mkconf

.PHONY: build
build: css | $(OUTPUT_DIR)
	go build -o $(OUTPUT_DIR)/$(SERVER_BIN) main.go

.PHONY: css
css:
	sass assets/stylesheets/main.scss:assets/stylesheets/styles.css

.PHONY: config
config: tools
	$(MKCONF_PATH) \
		--template configuration/anduril-config.json \
		--to $(OUTPUT_DIR)/data/encrypted-config.txt \
		--profile prod \
		--decrypt | jq

.PHONY: tools | $(OUTPUT_DIR)
tools: | $(OUTPUT_DIR)
	go build -o $(MKCONF_PATH) ./tools/mkconf.go

.PHONY: all
all: $(OUTPUT_DIR) build tools config

.PHONY: devenv
devenv: build tools
	mkdir -p $(OUTPUT_DIR)/data/assets
	rsync -rv assets/templates $(OUTPUT_DIR)/data/
	rsync -v assets/scripts/*.js $(OUTPUT_DIR)/data/assets/
	rsync -v assets/stylesheets/*.css $(OUTPUT_DIR)/data/assets/
	rsync -rv assets/icons $(OUTPUT_DIR)/data/assets/
	$(MKCONF_PATH) \
		--template configuration/anduril-config.json \
		--to $(OUTPUT_DIR)/data/encrypted-config.txt \
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
	mkdir -p $(OUTPUT_DIR)/data

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)/logs
	rm -rf $(OUTPUT_DIR)/work
	rm -rf $(OUTPUT_DIR)/data
	rm -f $(OUTPUT_DIR)/tls*
	rm -f $(OUTPUT_DIR)/$(SERVER_BIN)
	rm -f $(MKCONF_PATH)

.PHONY: tidy
tidy:
	go mod tidy
