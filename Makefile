# Makefile for Go module github.com/cicovic-andrija/anduril
#

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
	cp -r assets/templates $(OUTPUT_DIR)/data
	cp -r static $(OUTPUT_DIR)/data
	$(MKCONF_PATH) \
		--template configuration/anduril-config.json \
		--to $(OUTPUT_DIR)/data/encrypted-config.txt \
		--profile dev \
		--password v1.0.3-7841277 \
		--salt d3c24a79-c533-4e2d-974a-b4aab92198a6 \
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
