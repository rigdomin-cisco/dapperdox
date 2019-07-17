MAKEFLAGS += -r --warn-undefined-variables
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -euc
.DEFAULT_GOAL := help

include Makefile.variables

.PHONY: help clean veryclean build vendor dep-* format check test cover docs adhoc xcompile

## display this help message
help:
	@echo 'Management commands for dinamo:'
	@echo
	@echo 'Usage:'
	@echo '  ## Build Commands'
	@echo '    build           Compile the project.'
	@echo '    xcompile        Compile the project for multiple OS and Architectures.'
	@echo
	@echo '  ## Develop / Test Commands'
	@echo '    vendor          Install dependencies using dep if Gopkg.toml changed.'
	@echo '    dep-update      Update dependencies using dep.'
	@echo '    dep-add         Add new dependencies to dep and install.'
	@echo '    format          Run code formatter.'
	@echo '    check           Run static code analysis (lint).'
	@echo '    test            Run tests on project.'
	@echo '    cover           Run tests and capture code coverage metrics on project.'
	@echo '    clean           Clean the directory tree of produced artifacts.'
	@echo '    veryclean       Same as clean but also removes cached dependencies.'
	@echo

.ci-clean:
ifeq ($(CI_ENABLED),1)
	@rm -f tmp/dev_image_id || :
endif

## Clean the directory tree of produced artifacts.
clean: .ci-clean prepare
	@${DOCKERRUN} bash -c 'rm -rf bin build release cover vendor.orig *.out *.xml'

## Same as clean but also removes cached dependencies.
veryclean: clean
	@${DOCKERRUN} bash -c 'rm -rf tmp vendor'

## builds the dev container
prepare: tmp/dev_image_id
tmp/dev_image_id: Dockerfile.dev
	@mkdir -p tmp
	@docker rmi -f ${DEV_IMAGE} > /dev/null 2>&1 || true
	@echo "## Building dev container"
	@docker build --quiet -t ${DEV_IMAGE} -f Dockerfile.dev .
	@docker inspect -f "{{ .ID }}" ${DEV_IMAGE} > tmp/dev_image_id

# ----------------------------------------------
# build

## Compile the project.
build: build/dev

build/dev: check */*.go
	@rm -rf bin/
	@mkdir -p bin
	${DOCKERRUN} bash ./scripts/build.sh
	@chmod 755 bin/* || :

## Compile the project for multiple OS and Architectures.
xcompile: check
	@rm -rf build/
	@mkdir -p build
	${DOCKERRUN} bash ./scripts/xcompile.sh
	@find build -type d -exec chmod 755 {} \; || :
	@find build -type f -exec chmod 755 {} \; || :

# ----------------------------------------------
# dependencies

## Install dependencies using dep if Gopkg.toml changed.
vendor: tmp/vendor-installed
tmp/vendor-installed: tmp/dev_image_id Gopkg.toml
	@mkdir -p vendor
	${DOCKERRUN} dep ensure
	@date > tmp/vendor-installed
	@chmod 644 Gopkg.lock || :

## Update dependencies using dep.
dep-update: prepare
	${DOCKERRUN} dep ensure -update ${DEP}
	@chmod 644 Gopkg.lock || :

# usage DEP=github.com/owner/package make dep-add
## Add new dependencies to dep and install.
dep-add: prepare
ifeq ($(strip $(DEP)),)
	$(error "No dependency provided. Expected: DEP=<go import path>")
endif
	${DOCKERRUN} dep ensure -add ${DEP}
	@chmod 644 Gopkg.lock || :

# ----------------------------------------------
# develop and test

## print environment info about this dev environment
debug:
	@echo IMPORT_PATH="$(IMPORT_PATH)"
	@echo ROOT="$(ROOT)"
	@echo
	@echo docker commands run as:
	@echo "$(DOCKERRUN)"

## Run code formatter.
format: tmp/vendor-installed
	${DOCKERNOVENDOR} bash ./scripts/format.sh
	@if [[ -n "$$(git -c core.fileMode=false status --porcelain)" ]]; then \
    	echo "goimports modified code; requires attention!" ; \
    	if [[ "${CI_ENABLED}" == "1" ]]; then \
        	exit 1 ; \
    	fi ; \
	fi

## Run static code analysis (lint).
check: format
	${DOCKERNOVENDOR} bash ./scripts/check.sh

## Run tests on project.
test: check
	${DOCKERRUN} bash ./scripts/test.sh

## Run tests and capture code coverage metrics on project.
cover: check
	@rm -rf cover/
	@mkdir -p cover
ifeq ($(CI_ENABLED),1)
	${DOCKERRUN} bash ./scripts/cover.sh --ci
else
	${DOCKERRUN} bash ./scripts/cover.sh
	@chmod 644 cover/coverage.html
endif

# generate a TODO.md file with a list of TODO and FIXME items sorted by file
# the string is case insensitive and is removed from the output. So the final output
# should provide the file, line number, username that added it, and message about what
# needs to be done.
# Excludes the Makefile from consideration. Only files that are being tracked in git are
# included by default, therefore external dependencies or anything that is part of gitignore
# is automatically excluded.
## Generate a TODO list for project.
todo: prepare
	${DOCKERNOVENDOR} bash ./scripts/todo.sh -e Makefile -e scripts/todo.sh -e context.TODO -t '(FIXME|TODO)'

# usage: make adhoc RUNTHIS='command to run inside of dev container'
# example: make adhoc RUNTHIS='which jq'
adhoc: prepare
	@${DOCKERRUN} ${RUNTHIS}
