# .PHONY: build

# build:
# 	sam build

.PHONY: deps clean build

deps: 
	go get -u ./...

clean: 
	rm -rf ./showDocs/showDocs

build:
	GOOS=linux GOARCH=amd64 
	go build -o db/db ./db/main.go
