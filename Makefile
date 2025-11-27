export GO111MODULE=on
export CGO_ENABLED=0

APP_NAME := gmd

BUILD_DIR := dist
PLATFORMS := linux-amd64 linux-386 linux-arm linux-arm64

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
BUILDTIME := $(shell LANG=en_US.UTF-8; date +'%d %b %Y')

LOCAL_GOOS := $(shell go env GOOS)
LOCAL_GOARCH := $(shell go env GOARCH)

all: debug

release: $(LOCAL_GOOS)-$(LOCAL_GOARCH)

debug: $(BUILD_DIR)
	@EXT=""
	$(eval EXT=$(shell if [ "$(word 1,$(subst -, ,$@))" = "windows" ]; then echo .exe; fi))
	@echo "Building $@ with extension $(EXT)..."
	go build -o $(APP_NAME)$(EXT) -gcflags=all=-d=checkptr -ldflags '-X "github.com/kdruelle/gmd/cmd.version=${VERSION}" -X "github.com/kdruelle/gmd/cmd.buildDate=${BUILDTIME}"' ./

releases: clean
	@for t in $(PLATFORMS); do \
		$(MAKE) $$t; \
		tar zcvf $(BUILD_DIR)/$(APP_NAME)-$$t-$(VERSION).tar.gz -C $(BUILD_DIR) $(APP_NAME)-$$t-$(VERSION); \
	done

$(PLATFORMS): $(BUILD_DIR)
	@EXT=""
	$(eval EXT=$(shell if [ "$(word 1,$(subst -, ,$@))" = "windows" ]; then echo .exe; fi))
	@echo "Building $@ with extension $(EXT)..."
	GOOS=$(word 1,$(subst -, ,$@)) GOARCH=$(word 2,$(subst -, ,$@)) CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(APP_NAME)-$@-$(VERSION)/$(APP_NAME)$(EXT) -trimpath -ldflags '-s -w -X "cmd.version=${VERSION}" -X "cmd.buildTime=${BUILDTIME}"' ./

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

clean:
	@rm -rf $(BUILD_DIR)

version:
	@echo "Version: $(VERSION)"

lint: lint_install
	./bin/golangci-lint run

lint_install: bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s






ci:
	gox -osarch="linux/amd64 linux/386 linux/arm darwin/amd64 darwin/386 windows/amd64 windows/386" -output "./build/sslcap-${VERSION}-{{.OS}}-{{.Arch}}.bin" -ldflags '-s -w -X "cmd.version=${VERSION}" -X "cmd.buildTime=${BUILDTIME}"'
	@mkdir -p dist
	@export SSLCAP_VERSION=${VERSION}; ./release.sh



