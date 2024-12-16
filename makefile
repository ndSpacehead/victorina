.SILENT:

BINARY_NAME_UNIX = vic
BINARY_NAME_WIN = vic.exe

default: build

build:
	go build -v -ldflags '-s -w' \
		-o '.build/$(BINARY_NAME_UNIX)' .

build_linux:
	GOOS=linux GOARCH=amd64 go build -v -ldflags '-s -w' \
		-o '.build/$(BINARY_NAME_UNIX)' .

build_linux_arm64:
	GOOS=linux GOARCH=arm64 go build -v -ldflags '-s -w' \
		-o '.build/$(BINARY_NAME_UNIX)' .

build_win:
	GOOS=linux GOARCH=amd64 go build -v -ldflags '-s -w' \
		-o '.build/$(BINARY_NAME_WIN)' .

build_darwin:
	GOOS=darwin GOARCH=arm64 go build -v -ldflags '-s -w' \
		-o '.build/$(BINARY_NAME_UNIX)' .

build_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -v -ldflags '-s -w' \
		-o '.build/$(BINARY_NAME_UNIX)' .

run: build
	'.build/$(BINARY_NAME_UNIX)'