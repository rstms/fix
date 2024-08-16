#
# go makefile
#
#

build: fmt
	fix -- go build


fmt:
	fix -- go fmt ./...

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

README.md:
	@{ echo '# $(notdir $(PWD))'; echo '```';./$(notdir $(PWD)) --help; echo '```'; } >$@
