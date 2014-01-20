
server: server.go
	gofmt -w -tabs=false -tabwidth=4 server.go
	go build -o server server.go

clean:
	rm -f server