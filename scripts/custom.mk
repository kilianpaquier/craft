.PHONY: examples
examples: buildall
	@(cd examples/generic_github && ../../craft generate)
	@(cd examples/generic_gitlab && ../../craft generate)
	@(cd examples/golang_github && go mod tidy && ../../craft generate)
	@(cd examples/golang_gitlab && go mod tidy && ../../craft generate)
	@(cd examples/hugo_github && go mod tidy && ../../craft generate)
	@(cd examples/hugo_gitlab && go mod tidy && ../../craft generate)
	@(cd examples/helm && ../../craft generate)
	@(cd examples/nodejs_github && ../../craft generate)
	@(cd examples/nodejs_gitlab && ../../craft generate)