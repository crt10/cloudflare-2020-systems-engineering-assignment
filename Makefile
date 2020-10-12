build:
	go build .

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o .

clean:
	rm html-request