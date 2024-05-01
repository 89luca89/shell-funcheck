.PHONY: all shell-funcheck-amd64 shell-funcheck-arm64

all: shell-funcheck-amd64 shell-funcheck-arm64

shell-funcheck-amd64:
	@rm -f shell-funcheck-amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o shell-funcheck-amd64 main.go
	sha256sum shell-funcheck-amd64 > shell-funcheck-amd64.sha256

shell-funcheck-arm64:
	@rm -f shell-funcheck-arm64
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o shell-funcheck-arm64 main.go
	sha256sum shell-funcheck-arm64 > shell-funcheck-arm64.sha256
