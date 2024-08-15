#
# go makefile
#
modules := . ./vimfix

build: fmt
	fix -- go build


fmt:
	fix -- go fmt $(modules)

clean:
	go clean .

sterile: clean
	go clean -i -r -cache -testcache


test: build
	fix -- go test -v ./...

install:
	go install

uninstall:
	go clean -i
