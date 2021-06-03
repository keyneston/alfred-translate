.PHONY:
link: build
	go run github.com/jason0x43/go-alfred/alfred link

.PHONY:
build: $(shell find . -type f -name '*.go')
	go build -o workflow/alfred-translate
