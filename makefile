
# Force the use of bash
SHELL := bash

# Use a single shell per recipe instead of per line
.ONESHELL:

# -e Immediately exit if any command has a non-zero exit status
# -u Referencing any undefined variable is an error
# -o pipefail If any command in a pipeline fails, that return code will be used as the return code of the whole pipeline
.SHELLFLAGS := -eu -o pipefail -c

# If a recipe fails, delete the file
.DELETE_ON_ERROR:

MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

.PHONY: build
build:
	go build -C src -o ../bin/av-scanner

.PHONY: print-db
print-db:
	echo "SELECT * from files" | sqlite3 av-scanner-results.db

.PHONY: clear-db
clear-db:
	echo "DELETE from files" | sqlite3 av-scanner-results.db

.PHONY: scan
sweep:
	./bin/av-scanner scan --mode sweep --notify-endpoint http://ntfy.sh --notify-topic 1234-av-scanner-testing

.PHONY: watch
watch:
	./bin/av-scanner watch --mode watch --notify-endpoint http://ntfy.sh --notify-topic 1234-av-scanner-testing
