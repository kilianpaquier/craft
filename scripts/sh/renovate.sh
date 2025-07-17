#!/bin/sh

# renovate: datasource=golang-version
install-tool golang 1.24.5
make testdata
