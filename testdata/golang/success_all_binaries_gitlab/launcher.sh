#!/bin/sh
# Code generated by craft; DO NOT EDIT.

case $BINARY_NAME in
    cli-name) /app/cli-name;;
    cron-name) /app/cron-name;;
    job-name) /app/job-name;;
    worker-name) /app/worker-name;;
    *) echo "invalid binary '$BINARY_NAME'" && exit 1;;
esac