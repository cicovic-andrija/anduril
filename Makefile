# Makefile for Go module github.com/cicovic-andrija/anduril
#

VERSION = v1.0
BUILD = $(shell git rev-parse --short HEAD)
OUTPUT_DIR = out
SERVER_BIN = anduril-server

.PHONY: build
build: | $(OUTPUT_DIR)
	go build -v \
		-ldflags "-X github.com/cicovic-andrija/anduril/service/version=$(VERSION) -X github.com/cicovic-andrija/anduril/service/build=$(BUILD)" \
		-o $(OUTPUT_DIR)/$(SERVER_BIN) main.go

.PHONY: tools
tools: | $(OUTPUT_DIR)
	go build -o $(OUTPUT_DIR)/test-crypto ./tools/test-crypto.go

.PHONY: all
all: $(OUTPUT_DIR) build tools

.PHONY: devenv
devenv: build
	mkdir -p $(OUTPUT_DIR)/data
	cp -r templates $(OUTPUT_DIR)/data
	cp -r static $(OUTPUT_DIR)/data
	cp configuration/anduril-config-dev.json $(OUTPUT_DIR)/data/anduril-config.json
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
	rm -rf $(OUTPUT_DIR)

.PHONY: tidy
tidy:
	go mod tidy
