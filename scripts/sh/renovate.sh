#!/bin/sh

TESTDATA=1 go test ./... -count 1 -timeout=15s
