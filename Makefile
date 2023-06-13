
dev:
	go build -o baulos ./src/
	mv baulos /usr/local/bin

bin:
	- mkdir -p bin
	- GOOS=darwin GOARCH=amd64 go build -o bin/baulos_darwin_amd64
	- GOOS=darwin GOARCH=arm64 go build -o bin/baulos_darwin_arm64
	- GOOS=linux GOARCH=amd64 go build -o bin/baulos_linux_amd64
	- GOOS=linux GOARCH=arm64 go build -o bin/baulos_linux_arm64
