# Makefile for Go module github.com/cicovic-andrija/anduril
#

VERSION = v1.0
BUILD = $(shell git rev-parse --short HEAD)
OUT_DIR = out
SERVER_BIN = anduril-server

.PHONY: build
build: | $(OUT_DIR)
	@echo $(VERSION)-$(BUILD)
	go build -v -ldflags "-X github.com/cicovic-andrija/anduril/service/pass=$(VERSION)-$(BUILD)" -o $(OUT_DIR)/$(SERVER_BIN) main.go

.PHONY: confenc
confenc: | $(OUT_DIR)
	go build -ldflags "-X github.com/cicovic-andrija/anduril/service.pass=$(VERSION)-$(BUILD)" -o $(OUT_DIR)/$@ ./$@/...

$(OUT_DIR):
	mkdir -p $(OUT_DIR)

.PHONY: clean
clean:
	rm -rf $(OUT_DIR)

.PHONY: tidy
tidy:
	go mod tidy
