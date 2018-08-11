GO=go
LINTER=$(GOPATH)/bin/gometalinter
default: test

$(LINTER):
	$(GO) get -u github.com/alecthomas/gometalinter
	$(LINTER) --install &>/dev/null

lint: $(LINTER)
	$(LINTER)

test: lint
	$(GO) test -cover

.PHONY: test lint default
