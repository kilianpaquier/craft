.PHONY: examples
examples: buildall
	@(cd examples/generic_github && ../../craft generate --force-all)
	@(cd examples/generic_gitlab && ../../craft generate --force-all)
	@(cd examples/golang_github && go mod tidy && ../../craft generate --force-all)
	@(cd examples/golang_gitlab && go mod tidy && ../../craft generate --force-all)
	@(cd examples/hugo_github && go mod tidy && ../../craft generate --force-all)
	@(cd examples/hugo_gitlab && go mod tidy && ../../craft generate --force-all)
	@(cd examples/helm && ../../craft generate --force-all)
	@(cd examples/nodejs_github && ../../craft generate --force-all)
	@(cd examples/nodejs_gitlab && ../../craft generate --force-all)