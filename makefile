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
