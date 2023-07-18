
arm64: *.go go.mod
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags '-s -w -extldflags "-static"'

amd64: *.go go.mod
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w -extldflags "-static"'

upx:
	upx --lzma ipalf

clean:
	rm -f ipalf ipalf.upx

deps:
	go get -u -t ./...
