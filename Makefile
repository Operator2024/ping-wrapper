BINARY_NAME=ping_wrapper
VERSION=0.3.0
DATE=$(shell date -u +'%Y-%m-%d %I:%M:%S%p %Z')

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-macos -ldflags "-X 'main.version=${VERSION}' -X 'main.date=${DATE}'" main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-ubuntu -ldflags "-X 'main.version=${VERSION}' -X 'main.date=${DATE}'" main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-win.exe -ldflags "-X 'main.version=${VERSION}' -X 'main.date=${DATE}'" main.go

run:
	./${BINARY_NAME}-ubuntu

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}-macos
	rm ${BINARY_NAME}-ubuntu
	rm ${BINARY_NAME}-win.exe

