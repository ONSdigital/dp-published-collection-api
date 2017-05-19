
BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)
DATE:=$(shell date '+%Y%m%d-%H%M%S')
TGZ_FILE=dp-published-collection-api-$(GOOS)-$(GOARCH)-$(DATE)-$(HASH).tar.gz

build:
	@mkdir -p $(BUILD_ARCH)
	go build -o $(BUILD_ARCH)/bin/dp-published-collection-api cmd/dp-published-collection-api/main.go
	
test:
	go test -cover publishedcollection/* 
	
clean:
	rm -r	$(BUILD_ARCH) || true
	rm -r dp-published-collection-api-$(GOOS)-$(GOARCH)-*.tar.gz || true

package: build
	tar -zcf $(TGZ_FILE) -C $(BUILD_ARCH) .
	
hash:
	@git rev-parse --short HEAD
	
	
.PHONY: build clean