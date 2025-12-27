VERSION := 0.9.0

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


run: build
	./hanamark

init: 
	go mod init hanamark
	$(MAKE) run


# running files in hanamark and committing the code in repo 
# (add this make file in your markdown block or modify paths accordingly)
DATE := $(shell date '+%Y-%m-%d_%H-%M-%S')
makehanamark:
	cd .. && \
	cd hanamark && \
	make run && \
	cd .. && \
	cd thisisvoid/ && \
	git add . && \
	git commit -m "personal site_$(DATE)" &&\
	git push && \
	@echo "deployed hanamark parsed blog successfully..."

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
