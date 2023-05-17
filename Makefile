.PHONY: dev

STYLES = $(shell find frontend/styles -type f)
TEMPLS = $(shell find frontend -type f -name '*.templ')
GOFILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./build/*")

build: build/onlyserve

build/onlyserve: frontend internal/gh db $(GOFILES)
	mkdir -p build
	go build -o build/onlyserve ./cmd/onlyserve

build/github/schema.docs.graphql:
	mkdir -p build/github
	wget -P build/github/schema.docs.graphql -q -N https://docs.github.com/public/schema.docs.graphql


frontend: frontend/static/styles.css $(TEMPLS)
	cd frontend && templ generate
	touch frontend

frontend/static/styles.css: frontend/styles.scss $(STYLES)
	cd frontend && sass styles.scss static/styles.css


internal/gh: internal/gh/queries.graphql.gen.go

internal/gh/schema.docs.graphql.gz:
	mkdir -p internal/gh
	curl -sL https://docs.github.com/public/schema.docs.graphql | gzip - > internal/gh/schema.docs.graphql.gz

internal/gh/queries.graphql.gen.go: internal/gh/queries.graphql internal/gh/genqlient.yaml internal/gh/schema.docs.graphql.gz
	gzip -d internal/gh/schema.docs.graphql.gz
	cd internal/gh/ && genqlient
	gzip -f internal/gh/schema.docs.graphql


db: db/sqlitec

db/sqlitec: db/sqlc.json db/sqlitec/*.sql
	cd db && sqlc generate


dev:
	saq -- bash -c 'make && ./build/onlyserve --http localhost:8081'
