all:
	mkdir -p bin
	go build -o bin/standalone cmd/standalone/main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/lambda cmd/lambda/main.go
	rm -f bin/lambda.zip
	zip bin/lambda.zip bin/lambda