clean:
	rm -rf rdr.exe
build: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rdr-linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o rdr-linux-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o rdr-windows-amd64.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o rdr-windows-arm64.exe
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o rdr-darwin-amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o rdr-darwin-arm64
