all: fmt build package push

fmt:
	gofmt -w *.go cmd/*.go

build:
	env GOOS=linux GOARCH=amd64 go build -o watch cmd/main.go

package:
	rm -rf alert
	mkdir -p alert
	cp watch alert/
	cp config.json alert/
	cp smtp.json alert/
	tar czf alert.tar.gz alert
