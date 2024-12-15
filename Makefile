# Change these variables as necessary.
APP=website
APP_EXE="./out/$(APP)"

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy

## Build: build the go application
.PHONY: build
build:
	mkdir -p out/
	go build -o $(APP_EXE)
	@echo "Build passed"

## Run: runs the go binary. use additional options if required.
run:
	make build
	$(APP_EXE)
