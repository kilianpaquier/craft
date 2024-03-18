# Code generated by craft; DO NOT EDIT.

GCI_CONFIG_PATH := .golangci.yml

.PHONY: reports
reports:
	@mkdir -p reports/

.PHONY: lint
lint: reports
	@golangci-lint run -c ${GCI_CONFIG_PATH} --timeout 240s --fast --sort-results \
		--out-format checkstyle:reports/go-ci-lint.checkstyle.xml,colored-line-number $(ARGS) || \
		echo "golangci-lint failed, running 'make lint-fix' may fix some issues"

.PHONY: lint-fix
lint-fix: reports
	@ARGS="--fix" make -s lint

.PHONY: test
test: lint
	@go test ./... -count 1

.PHONY: test-race
test-race: lint
	@go test ./... -race

.PHONY: test-cover
test-cover: lint reports
	@go test ./... -coverpkg="./..." -covermode="count" -coverprofile="reports/go-coverage.native.out"

.PHONY: buildall
buildall: build-cron-refresh build-worker-appli

.PHONY: cron-refresh worker-appli
build-%:
	@CGO_ENABLED=0 go build -o $* cmd/$*/main.go

.PHONY: cron-refresh worker-appli
local-%:
	@go run cmd/$*/main.go

build-docker:
	@docker build -t golang-app .