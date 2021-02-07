GOOS ?= linux
GOARCH ?= arm

.PHONY: all clean

all: revo-web-dispatch

revo-web-dispatch: revo-web-dispatch.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-s -w" revo-web-dispatch.go
	upx --brute revo-web-dispatch

clean:
	rm -f revo-web-dispatch
