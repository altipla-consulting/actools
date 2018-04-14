
FILES = $(shell find . -type f -name "*.go" -not -path "./vendor/*")

gofmt:
	@gofmt -w $(FILES)
	@gofmt -r '&a{} -> new(a)' -w $(FILES)
