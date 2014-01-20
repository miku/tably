README
======

A webserver that provides a rough deduplication dashboard.

Develop
-------

	$ git clone git@github.com:miku/tably.git
	$ cd tably
	$ # make sure you have a GOPATH (e.g. $HOME or $HOME/go)
	$ make
	$ ./server
	[martini] listening on port 3000
	...

Adjust `example.json` configuration and rename it to `server.json`.

Go to http://localhost:3000 --