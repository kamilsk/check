WD = $(shell if ! [ -e main.go ]; then echo $(dir $(PWD)); else echo $(PWD); fi)

.PHONY: cmd-urls
cmd-urls:
	cd $(WD) && go run main.go build.go urls
