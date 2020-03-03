
FILES = $(shell find . -type f -name "*.go" -not -path "./vendor/*")

gofmt:
	@gofmt -s -w $(FILES)
	@gofmt -r '&a{} -> new(a)' -w $(FILES)
