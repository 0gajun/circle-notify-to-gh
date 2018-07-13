bin_name=circleci-comments-to-pr

.PHONY: build
build:
	go build -o ${bin_name} .

.PHONY: build-for-linux
build-for-linux:
	GOOS=linux GOARCH=amd64 go build -o ${bin_name} .
