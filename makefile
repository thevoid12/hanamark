build:
	echo "building..."
	go mod tidy
	go build -o hanamark

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
