.PHONY: help test test-race test-leak bench bench-compare bench-swagger-gen lint sec-scan upgrade release release-tag changelog-gen changelog-commit

help: ## show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

PROJECT_NAME?=fiat
APP_NAME?=bnpl

########
# test #
########

test: test-race test-leak ## launch all tests 

test-race: ## launch all tests with race detection
	go test ./... -cover -race

test-leak: ## launch all tests with leak detection (if possible)
	go test ./internal/port/rest/userfacing/... -leak
	go test ./internal/port/rest/internalfacing/... -leak

#############
# benchmark #
#############

bench: ## launch benchs
	go test ./... -bench=. -benchmem | tee ./bench.txt

bench-compare: ## compare benchs results
	benchstat ./bench.txt

bench-swagger-gen: ## generate code to benchmark your deployed api with k6 (see https://k6.io/blog/load-testing-your-api-with-swagger-openapi-and-k6/)
	@printf "what is the url of the json swagger api definition? (for example http://localhost:3000/v1/swagger/doc.json) "; read SWAGGER_URL &&\
	openapi-generator generate \
		-i $$SWAGGER_URL \
		-g k6 \
		-o ./benchmarks/ &&\
	printf "you can now run the following command: k6 run -d 5s -u 100 ./benchmarks/script.js"

########
# lint #
########

lint: ## lints the entire codebase
	@golangci-lint run ./... --config=./.golangci.toml && \
	if [ $$(gofumpt -e -l ./ | wc -l) == "0" ] ; \
		then exit 0; \
	else \
		echo "these files needs to be gofumpt-ed"; \
		gofumpt -e -l ./; \
		exit 1; \
	fi

#######
# sec #
#######

sec-scan: ## scan for sec issues with trivy (trivy binary needed)
	trivy fs --exit-code 1 --no-progress --severity CRITICAL ./

############
# upgrades #
############

upgrade: ## upgrade dependencies (beware, it can break everything)
	go mod tidy && \
	go get -t -u ./... && \
	go mod tidy

#########
# build #
#########

build: lint test bench sec-scan docker-build ## lint, test, bench and sec scan before building the docker image
	@printf "\nyou can now deploy to your env of choice:\ncd deploy\nENV=dev make deploy-latest\n"

LAST_MAIN_COMMIT_HASH=$(shell git rev-parse --short HEAD)
LAST_MAIN_COMMIT_TIME=$(shell git log main -n1 --format='%cd' --date='iso-strict')

docker-build-ci: ## docker build, works only in the cloud
	DOCKER_BUILDKIT=1 \
	docker build \
		-f Dockerfile \
		-t $(APP_NAME) \
		--build-arg LAST_MAIN_COMMIT_HASH=$(LAST_MAIN_COMMIT_HASH) \
		--build-arg LAST_MAIN_COMMIT_TIME=$(LAST_MAIN_COMMIT_TIME) \
		--ssh default \
		--progress=plain \
		./

docker-build-local: ## docker build locally, works on m1 macs
	@printf "what is your github username: "; read -r GITHUB_USER && \
	printf "what is your github personal access token: "; read -rs GITHUB_PERSONAL_ACCESS_TOKEN && \
	docker build \
		-f Dockerfile.local \
		-t $(APP_NAME) \
		--build-arg LAST_MAIN_COMMIT_HASH=$(LAST_MAIN_COMMIT_HASH) \
		--build-arg LAST_MAIN_COMMIT_TIME=$(LAST_MAIN_COMMIT_TIME) \
		--build-arg GITHUB_USER=$$GITHUB_USER \
		--build-arg GITHUB_PERSONAL_ACCESS_TOKEN=$$GITHUB_PERSONAL_ACCESS_TOKEN \
		./

###########
# release #
###########

release: release-tag changelog-gen changelog-commit ## create a new tag to release this module

MOD_VERSION = $(shell git describe --abbrev=0 --tags")
	
release-tag: 
	@printf "here is the latest tag present: "; \
	printf "$(MOD_VERSION)\n"; \
	printf "what tag do you want to give? (use the form vX.X.X): "; \
	read -r TAG && \
	git tag $$TAG && \
	printf "\nrelease tagged $$TAG !\n"

#############
# changelog #
#############

MESSAGE_CHANGELOG_COMMIT="update CHANGELOG.md for $(MOD_VERSION)"

changelog-gen: ## generates the changelog in CHANGELOG.md
	@git cliff \
		-o ./CHANGELOG.md && \
	printf "\nchangelog generated!\n"

# keep this commit unconventional so it doesnt appear in the changelog
changelog-commit:
	git commit -m $(MESSAGE_CHANGELOG_COMMIT) ./CHANGELOG.md

######
# db #
######

db-pg-init: 
	@( \
	printf "Enter pass for db: "; read -rs DB_PASSWORD && \
	printf "\nEnter environment suffix(_dev, _local...): "; read DB_SUFFIX &&\
	sed \
	-e "s/DB_PASSWORD/$$DB_PASSWORD/g" \
	-e "s/DB_SUFFIX/$$DB_SUFFIX/g" \
	./db/init/init.sql | \
	PGPASSWORD=$$DB_PASSWORD psql -h localhost -p 5436 -U postgres -f - \
	)

db-cockroachdb-rootkey:
	mkdir ./db/crdb-certs && \
	kubectl cp cockroachdb/cockroachdb-0:cockroach-certs/ca.crt ./db/crdb-certs/ca.crt -c db && \
	cockroach cert create-client \
		--certs-dir=./db/crdb-certs \
		--ca-key=$(CAROOT)/rootCA-key.pem root

db-cockroachdb-init:
	@( \
	printf "Enter pass for db: "; read -s DB_PASSWORD && \
	printf "\nEnter environment suffix(_dev, _local...): "; read DB_SUFFIX &&\
	printf "Enter port(26257...): "; read -r DB_PORT &&\
	sed \
	-e "s/DB_PASSWORD/$$DB_PASSWORD/g" \
	-e "s/DB_SUFFIX/$$DB_SUFFIX/g" \
	./db/init/init.sql > ./db/crdb-certs/init.sed.sql && \
	cockroach sql --certs-dir=./db/crdb-certs -f ./db/crdb-certs/init.sed.sql -p $$DB_PORT && \
	rm ./db/crdb-certs/init.sed.sql \
	)

#######
# sql #
#######

sql-gen-sqlboiler:
	@( \
	printf "Enter pass for db: "; read -s DB_PASSWORD && \
	PSQL_PASS=$$DB_PASSWORD sqlboiler psql -c ./db/sqlboiler.toml \
	)

###########
# swagger #
###########

swagger-gen:
	swag init -d ./cmd/api --parseDependency --outputTypes go --output "./internal/docs/"
