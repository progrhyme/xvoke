# Targets to build, check, and package the xvoke command.
# The release workflow (.github/workflows/release.yml) calls the same targets.
#
# This Makefile uses only what GNU make 3.81, which macOS has, supports.

GO ?= go
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)

# The version of a release, without the leading "v", such as 0.0.3. The archive targets need it.
VERSION ?=

NAME := xvoke
DIST := dist
PACKAGE := $(NAME)_$(VERSION)_$(GOOS)_$(GOARCH)

ifeq ($(GOOS),windows)
EXE := .exe
else
EXE :=
endif

SHA256 := $(shell command -v sha256sum >/dev/null 2>&1 && echo sha256sum || echo "shasum -a 256")

.PHONY: build test lint archive check-binary checksums require-version

# Build the command for this machine into dist/.
build:
	$(GO) build -o $(DIST)/$(NAME) ./cmd/$(NAME)

# Run go vet and go test.
test:
	$(GO) vet ./...
	$(GO) test ./...

lint:
	golangci-lint run

# Cross-compile for GOOS and GOARCH, and make dist/xvoke_VERSION_GOOS_GOARCH.tar.gz (.zip for Windows).
# The archive has the executable and LICENSE.
# The version of the command is the module version that Go records at build time, so build at a tagged commit.
archive: require-version
	mkdir -p $(DIST)/$(PACKAGE)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -trimpath -o $(DIST)/$(PACKAGE)/$(NAME)$(EXE) ./cmd/$(NAME)
	cp LICENSE $(DIST)/$(PACKAGE)/
ifeq ($(GOOS),windows)
	cd $(DIST) && zip -qr $(PACKAGE).zip $(PACKAGE)
else
	tar -czf $(DIST)/$(PACKAGE).tar.gz -C $(DIST) $(PACKAGE)
endif

# Check that the built command reports VERSION. It runs the binary, so GOOS and GOARCH must match this machine.
check-binary: require-version
	@reported="$$($(DIST)/$(PACKAGE)/$(NAME)$(EXE) version | jq -r .version)"; \
	if [ "$$reported" != "$(VERSION)" ]; then \
		echo "$(NAME) version reports $$reported, want $(VERSION)" >&2; exit 1; \
	fi

# Make dist/xvoke_VERSION_checksums.txt, which lists the SHA-256 of every archive of VERSION in dist/.
checksums: require-version
	@archives="$(notdir $(wildcard $(DIST)/$(NAME)_$(VERSION)_*.tar.gz $(DIST)/$(NAME)_$(VERSION)_*.zip))"; \
	if [ -z "$$archives" ]; then echo "no archive of $(VERSION) in $(DIST)/" >&2; exit 1; fi; \
	cd $(DIST) && $(SHA256) $$archives > $(NAME)_$(VERSION)_checksums.txt

require-version:
	@test -n "$(VERSION)" || { echo "VERSION is required, such as VERSION=0.0.3" >&2; exit 1; }
