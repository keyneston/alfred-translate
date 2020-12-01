.PHONY:
link: build
	alfred link

.PHONY:
build: $(shell find . -type f -name '*.go')
	go build -o workflow/alfred-translate
