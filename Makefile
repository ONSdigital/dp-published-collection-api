
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
	
nomad:
	@for t in *-template.nomad; do			\
	plan=$${t%-template.nomad}.nomad;	\
	test -f $$plan && rm $$plan;		\
	sed	-e 's,DATA_CENTER,$(DATA_CENTER),g'	        	\
		-e 's,S3_TAR_FILE,$(S3_TAR_FILE),g'		        	\
		-e 's,PUBLISH_DATABASE_URL,$(DATABASE_URL),g'		\
		< $$t > $$plan || exit 2;			\
done
	
	
.PHONY: build clean