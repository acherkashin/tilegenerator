test:
	go test ./...

install:
	go install ./...

install-deps:
	go get ./...

build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o tilegenerator.exe app.go
build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o tilegenerator app.go
build-all: build-windows build-linux
