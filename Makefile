# HELP
.PHONY: help

help: ## Usage: make <option>
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


build: ## APP Build.
	go build -o bin/search  cmd/search.go;
# 	docker-compose run --rm golang go build -o bin/search cmd/search.go

test: ## Test.
	go test -v -race -cover -coverprofile=var/log/coverage-search.out ./search/;
	go test -v -race -cover -coverprofile=var/log/coverage-http.out ./http/;
# 	docker-compose run --rm golang go test -v -race -cover -coverprofile=var/log/coverage-search.out ./search/;
# 	docker-compose run --rm golang go test -v -race -cover -coverprofile=var/log/coverage-http.out ./http/;

cover: ## Test coverage.
	go tool cover -func=var/log/coverage-search.out;
	go tool cover -func=var/log/coverage-http.out;
# 	docker-compose run --rm golang go tool cover -func=var/log/coverage-search.out;
# 	docker-compose run --rm golang go tool cover -func=var/log/coverage-http.out;

cover-html: ## Test coverage HTML.
	go tool cover -html=var/log/coverage-search.out
	go tool cover -html=var/log/coverage-http.out
# 	docker-compose run --rm golang go tool cover -html=var/log/coverage-search.out
# 	docker-compose run --rm golang go tool cover -html=var/log/coverage-http.out
