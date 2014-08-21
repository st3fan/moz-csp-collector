# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/

app: moz-csp-collector

moz-csp-collector: main.go csp/database.go csp/reports.go csp/server.go
	go build .

docker-image: moz-csp-collector
	docker build -t st3fan/moz-csp-collector .

docker-run:
	docker run -i -t -p "8080:8080" st3fan/moz-csp-collector

all: app docker-image

clean:
	rm -f moz-csp-collector

install_deps:
	go get -u github.com/lib/pq
	go get -u github.com/go-martini/martini
	go get -u github.com/martini-contrib/binding
	go get -u github.com/martini-contrib/render
