build:
	echo "building..."
	go mod tidy
	go build -o hanamark

run: build
	./hanamark

init: 
	go mod init hanamark
	$(MAKE) run
