VERSION = 0.0.1

.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

bin/github-to-terraform: ## Build github-to-terraform
	go build -ldflags "$(LDFLAGS) -X main.version=${VERSION}" -o bin/github-to-terraform .

.PHONY: clean
clean: ## Clean all generated files
	rm -rf ./bin
	rm -rf ./dist

.PHONY: goreleaser_check
goreleaser_check: ## Goreleaser check configuration
	goreleaser check

.PHONY: goreleaser_build
goreleaser_build: ## Goreleaser build configuration
	goreleaser build --snapshot --rm-dist

.PHONY: goreleaser_test
goreleaser_test: ## Goreleaser test configuration
	goreleaser --snapshot --skip-publish --rm-dist

default: help
