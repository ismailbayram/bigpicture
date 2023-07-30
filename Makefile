build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/bigpicture-amd64-darwin main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/bigpicture-amd64-linux main.go
