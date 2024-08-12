GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOFMT := "goimports"

fmt: ## Run gofmt for all .go files
	$(GOFMT) -w $(GOFMT_FILES)

test: ## Run go test for whole project
	go test -v ./...

create-migrate: ## Create a migration file
	migrate create -ext sql -dir db/migrations -seq $(filename)

migrateup: ## Migrate database up. Sample: make migrateup db="postgresql://[<username>][:<password>]@[<host>:<port>]/<dbname>?sslmode=disable"
	migrate -database ${db} -path db/migrations -verbose up ${N}

migratedown: ## Migrate database down. Sample: make migratedown db="postgresql://[<username>][:<password>]@[<host>:<port>]/<dbname>?sslmode=disable"
	migrate -database ${db} -path db/migrations -verbose down ${N}

goupdate: ## Remove go.sum and Go get dependencies
	@GOSUMDB=off rm -rf go.sum
	@GOSUMDB=off go mod tidy

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'