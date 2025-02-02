#
# go makefile
#
#

GO = fix -- go
#GO = go

build: fmt
	$(GO) build


fmt:
	$(GO) fmt ./...

clean:
	go clean .

sterile: clean
	go clean -i -r -cache -testcache


test: build
	$(GO) test -v ./...

install:
	go install

uninstall:
	go clean -i

README.md:
	@{ echo '# $(notdir $(PWD))'; echo '```';./$(notdir $(PWD)) --help; echo '```'; } >$@
