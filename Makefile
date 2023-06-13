
dev-arm:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags '-extldflags "-static"' -trimpath -o bin/zen_darwin_arm64 ./src/zen.go
	mv bin/zen_darwin_arm64 /usr/local/bin/zen

build:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o bin/zen_darwin_x86_64 ./src/zen.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags '-extldflags "-static"' -o bin/zen_darwin_arm64 ./src/zen.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o bin/zen_linux_x86_64 ./src/zen.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags '-extldflags "-static"' -o bin/zen_linux_arm64 ./src/zen.go