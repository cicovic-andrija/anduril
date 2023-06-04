# Makefile for Go module github.com/cicovic-andrija/anduril
#

VERSION = v1.0.0-b6c0821
BUILD = 94364bb8-7f08-43c5-aaca-42ec17ea18b4
OUTPUT_DIR = out
SERVER_BIN = anduril-server
MKCONF_PATH = $(OUTPUT_DIR)/mkconf

.PHONY: build
build: | $(OUTPUT_DIR)
	go build -o $(OUTPUT_DIR)/$(SERVER_BIN) main.go

.PHONY: tools
tools: | $(OUTPUT_DIR)
	go build -o $(MKCONF_PATH) ./tools/mkconf.go

.PHONY: all
all: $(OUTPUT_DIR) build tools

.PHONY: devenv
devenv: build tools
	mkdir -p $(OUTPUT_DIR)/data
	cp -r templates $(OUTPUT_DIR)/data
	cp -r static $(OUTPUT_DIR)/data
	$(MKCONF_PATH) \
		--template configuration/anduril-config.json \
		--to $(OUTPUT_DIR)/data/encrypted-config.txt \
		--profile dev \
		--password $(VERSION) \
		--salt $(BUILD) \
		--decrypt
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
	mkdir -p $(OUTPUT_DIR)

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
