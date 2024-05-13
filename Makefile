# Variables
MAIN_PACKAGE_PATH := ./cmd/api
BINARY_NAME := halo-suster
PWD := $(shell pwd)
tag := latest
DB_TYPE := postgres
DB_USERNAME := postgres
DB_PASSWORD := password
DB_HOST := localhost
DB_PORT := 5432
DB_NAME := halo-suster
DB_PARAMS := sslmode=disable


# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: run formatting, go mod tidy and goimports
.PHONY: tidy
tidy:
	go run mvdan.cc/gofumpt@latest -extra -l -w .
	go run github.com/segmentio/golines@latest --max-len=120 --shorten-comments -w .
	go run github.com/incu6us/goimports-reviser/v3@latest -rm-unused ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	#go run github.com/roblaszczak/go-cleanarch@latest -domain domain -application service -infrastructure repository -interfaces handler
	go run go.uber.org/nilaway/cmd/nilaway@latest ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run
	go test -race -buildvcs -vet=off ./...

## update: update all dependencies
.PHONY: update
update:
	go get -u ./...
	go mod tidy
	go mod verify


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs -count=1 ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -count=1 -coverprofile=./tmp/coverage.out ./...
	go tool cover -html=./tmp/coverage.out

## build: build the application
.PHONY: build
build:
	# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -o=./tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the  application
.PHONY: run
run: build
	./tmp/bin/${BINARY_NAME}

## watch: run the application with reloading on file changes
.PHONY: watch
watch:
	go run github.com/cosmtrek/air@latest \
		--build.cmd "make build" --build.bin "./tmp/bin/${BINARY_NAME}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, mod, tpl, tmpl, html, env" \
		--build.send_interrupt "true" \
		--build.kill_delay "5000000" \
		--misc.clean_on_exit "true"

## clean: remove the binary
.PHONY: clean
clean:
	rm -f ./tmp/bin/${BINARY_NAME}


# ==================================================================================== #
# DATABASE
# ==================================================================================== #

## db/connect: run the local database from docker compose
.PHONY: db/connect
db/connect:
	docker compose up -d postgres-project-sprint-2

## migrate/create name=$1: create a new migration
.PHONY: migrate/create
migrate/create:
	@echo "---Creating migration files---"
	go run -tags '$(DB_TYPE)' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
 		create -ext sql -dir $(PWD)/migrations -seq -digits 5 $(name)

## migrate/up n=$1: run the up migration
.PHONY: migrate/up
migrate/up:
	go run -tags '$(DB_TYPE)' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
 		-database '$(DB_TYPE)://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)' -path $(PWD)/migrations up $(n)

## migrate/down n=$1: run the down migration
.PHONY: migrate/down
migrate/down: confirm
	go run -tags '$(DB_TYPE)' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
 		-database '$(DB_TYPE)://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)' -path $(PWD)/migrations down $(n)

## migrate/drop: drop all tables
.PHONY: migrate/drop
migrate/drop: confirm
	go run -tags '$(DB_TYPE)' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
 		-database '$(DB_TYPE)://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)' -path $(PWD)/migrations drop

## migrate/goto version=$1: migrate to a specific version number
.PHONY: migrate/goto
migrate/goto: confirm
	go run -tags '$(DB_TYPE)' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
 		-database '$(DB_TYPE)://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)' -path $(PWD)/migrations goto $(version)

## migrate/force version=$1: force a migration by version
.PHONY: migrate/force
migrate/force: confirm
	go run -tags '$(DB_TYPE)' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
 		-database '$(DB_TYPE)://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)' -path $(PWD)/migrations force $(version)

## migrate/version: print the current migration version
.PHONY: migrate/version
migrate/version:
	go run -tags '$(DB_TYPE)' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
	 	-database '$(DB_TYPE)://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)' -path $(PWD)/migrations version


# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## production/build: build the application for production
.PHONY: production/build
production/build: confirm tidy audit
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o=./bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## docker/build tag=$1: build the application for docker
.PHONY: docker/build
docker/build: confirm update tidy audit
	docker build -t j03hanafi/${BINARY_NAME}:${tag} .

## docker/run tag=$1: run the application in a docker container
.PHONY: docker/run
docker/run:
	docker run -p 8080:8080 -d j03hanafi/${BINARY_NAME}:${tag}

## docker/push tag=$1: push the docker image to the registry
.PHONY: docker/push
docker/push: confirm
	docker push j03hanafi/${BINARY_NAME}:${tag}