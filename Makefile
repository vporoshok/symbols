.PHONY: help test cover

all: test cover

help: ## This message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run tests
	CGO_ENABLED=1 go test -race -cover -count=1 -covermode=atomic -coverprofile=coverage.out ./...

cover: test ## Check coverage
	@sed -i '/\(pb\|easyjson\|string\)\.go/d' coverage.out
	@go tool cover -func=.coverprofile | tail -n 1 | awk '{print "Total coverage:", $$3;}'
	@test `go tool cover -func=.coverprofile | tail -n 1 | awk '{print $$3;}' | sed 's/\..*//'` -ge 70
