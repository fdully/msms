run:
	go run cmd/sms-server/main.go

test:
	go test -v -race ./... -count=1

build:
   GOOS=linux GOARCH=amd64 go build -o smsgwapp cmd/sms-server/main.go

