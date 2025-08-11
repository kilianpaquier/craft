#!/bin/sh

# renovate: datasource=golang-version
install-tool golang 1.24.6
make testdata
