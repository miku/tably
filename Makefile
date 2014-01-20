server: server.go deps.lock
	gofmt -w -tabs=false -tabwidth=4 server.go
	go build -o server server.go

deps.lock:
	go get github.com/codegangsta/martini
	go get github.com/codegangsta/martini-contrib/render
	go get github.com/go-sql-driver/mysql
	go get github.com/jmoiron/sqlx
	touch deps.lock # remove this to redownload deps

clean-deps:
	rm -f deps.lock

clean:
	rm -f server