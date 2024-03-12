# This file is safe to edit. Once it exists it will not be overwritten.

include ./scripts/*.mk

.PHONY: examples
examples:
	@(cd examples/generic && go run ../../cmd/craft/main.go generate --force-all)
	@(cd examples/generic-helm && go run ../../cmd/craft/main.go generate --force-all)
	@(cd examples/golang-api && go run ../../cmd/craft/main.go generate --force-all)
	@(cd examples/golang-app && go run ../../cmd/craft/main.go generate --force-all)