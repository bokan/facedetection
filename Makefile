lint:
	golangci-lint run -E golint -E gofmt -E govet -E goimports -E gocyclo -E gocognit -E ineffassign -E misspell

