# Change these variables as necessary.
main_package_path = ./main.go
binary_name = mongo-go
BIN_LOCATION := /tmp

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message (DEFAULT)
.PHONY: help
help:
	@echo -e "Make commands for ${binary_name}\n"
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /' | sort

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

clean: confirm
	rm -rf ./tmp

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	@echo "Tidying up..."
	@go mod tidy -v
	@go fmt ./...

## build: build the application
.PHONY: build
build:
	@echo "Building ${binary_name}"
	@# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here... | for some reason, it is printed on the terminal so made it a terminal command to remove it
	@go build -o=${BIN_LOCATION}/${binary_name} ${main_package_path}
	@echo "Built ${binary_name}"

## run: run the  application
.PHONY: run
run: build
	@echo "Running ${binary_name}"
	@${BIN_LOCATION}/${binary_name} 

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	@echo "Running live reload using air..."
	@air \
		--build.cmd "make build" --build.bin "${BIN_LOCATION}/${binary_name}" \
		--build.exclude_dir "data" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true" \
		--screen.clear_on_rebuild "true"


# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: confirm audit no-dirty
	git push

## production/deploy: deploy the application to production
.PHONY: production/deploy
production/deploy: confirm audit no-dirty
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=/build/linux_amd64/${binary_name} ${main_package_path}
	# upx -5 /tmp/bin/linux_amd64/${binary_name} 
	# Include additional deployment steps here...

