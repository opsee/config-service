default: install

deps:
	@go get github.com/tools/godep

bindir:
	@mkdir -p bin/

build: bindir deps
	go build -o bin/cf .
	GOOS=linux GOARCH=amd64 go build -o bin/cf-linux-amd64 .

install: build
	@cp bin/* ${GOPATH}/bin

release: build
	@aws s3 cp bin/cf-linux-amd64 s3://opsee-releases/go/cf/cf-linux-amd64

clean:
	@rm -f bin/*
