PKG_DIR  = $(shell if ! [ -e main.go ]; then echo $(dir $(PWD)); else echo $(PWD); fi)
GO_FILES = main.go build.go

.PHONY: cmd-help
cmd-help:
	cd $(PKG_DIR) && go run $(GO_FILES) help

.PHONY: cmd-urls
cmd-urls:
	cd $(PKG_DIR) && go run $(GO_FILES) urls https://kamil.samigullin.info/

.PHONY: cmd-urls-help
cmd-urls-help:
	cd $(PKG_DIR) && go run $(GO_FILES) urls --help
