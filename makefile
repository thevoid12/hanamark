VERSION := $(shell grep -oE '[0-9]+\.[0-9]+\.[0-9]+' constant/version.go)

ifneq ("$(wildcard .env)","")
    include .env
    export
endif

build:
	echo "building executable version $(VERSION)..."
	go mod tidy
	@echo "package constants\n\nconst Version = \"$(VERSION)\"" > constant/version.go
	# build for linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hanamark
# 	# build for intel mac
# 	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o hanamark
# 	# build for apple silicon mac
# 	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o hanamark
# 	# build for windows
# 	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o hanamark.exe

# ---- Config ----
GOLANGCI_LINT := golangci-lint
GO := go

# ---- Targets ----
.PHONY: help lint lint-fix lint-install test vet tidy check

help:
	@echo "Available targets:"
	@echo "  make lint        Run golangci-lint"
	@echo "  make lint-fix    Run golangci-lint with autofix"
	@echo "  make test        Run go tests"
	@echo "  make vet         Run go vet"
	@echo "  make tidy        Run go mod tidy"
	@echo "  make check       Run all checks"

lint:
	$(GOLANGCI_LINT) run

lint-fix:
	$(GOLANGCI_LINT) run --fix

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

check: tidy vet test lint



# Makefile

.PHONY: release-check
release-check:
	@echo "Releasing version: $(VERSION)"

.PHONY: build-dist
build-dist:
	# This creates the binaries and changelog in ./dist/
	# --skip=publish ensures nothing is uploaded automatically
	goreleaser release --clean --snapshot --skip=publish

.PHONY: push-codeberg
push-codeberg:
	# Ensure codeberg remote is added: git remote add codeberg <url>
	git push codeberg main
	git push codeberg v$(VERSION)

.PHONY: full-release
full-release: release-check build-dist push-codeberg
	# 1. Run goreleaser for GitHub (using your .goreleaser.yaml)
	goreleaser release --clean
	# 2. Reminder for Codeberg
	@echo "Binary generation complete. Upload files from ./dist/ to Codeberg manually or via API."

.PHONY: add-tag 
add-tag:
	git tag -a v$(VERSION) -m "v$(VERSION)"
	git push origin v$(VERSION)

.PHONY: full-release-draft
full-release-draft:add-tag release-check build-dist 
	# 1. Run goreleaser for GitHub (using your .goreleaser.yaml)
	goreleaser release --clean --draft
	@echo "manually check the changes in github and click publish"
	# 2. Reminder for Codeberg
	@echo "Binary generation complete. Upload files from ./dist/ to Codeberg manually or via API."


