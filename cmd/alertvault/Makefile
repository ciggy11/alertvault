COMMIT_ID := `git log -1 --format=%H`
COMMIT_DATE := `git log -1 --format=%aI`
VERSION := $${CI_COMMIT_TAG:-SNAPSHOT-$(COMMIT_ID)}
SHELL := /bin/bash

GOOS ?= linux
GOARCH ?= amd64

.PHONY: help
help:	### help help help
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@# linux only xD
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

.PHONY: build
build: ### build the binary applying the correct version from git
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o alertvault-$(GOARCH)

CURRENT_DIR:=$(shell pwd)