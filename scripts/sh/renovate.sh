#!/bin/sh

TESTDATA=1 /usr/local/bin/go test ./... -count 1 -timeout=15s
