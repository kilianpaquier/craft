# This file is safe to edit. Once it exists it will not be overwritten.

include ./scripts/*.mk

.PHONY: examples
examples:
	@(cd examples/generic_github && go run ../../cmd/craft/main.go generate --force-all)
	@(cd examples/generic_gitlab && go run ../../cmd/craft/main.go generate --force-all)
	@(cd examples/golang_github && go mod tidy && go run ../../cmd/craft/main.go generate --force-all)
	@(cd examples/golang_gitlab && go mod tidy && go run ../../cmd/craft/main.go generate --force-all)
	@(cd examples/hugo_github && go mod tidy && go run ../../cmd/craft/main.go generate --force-all)
	@(cd examples/hugo_gitlab && go mod tidy && go run ../../cmd/craft/main.go generate --force-all)
	@(cd examples/helm && go run ../../cmd/craft/main.go generate --force-all)
	@(cd examples/nodejs_github && go run ../../cmd/craft/main.go generate --force-all)
	@(cd examples/nodejs_gitlab && go run ../../cmd/craft/main.go generate --force-all)