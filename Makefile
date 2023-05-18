.PHONY: build dev

STYLES = $(shell find frontend/styles -type f)
TEMPLS = $(shell find frontend -type f -name '*.templ')
SCRIPTS = $(shell find frontend/scripts -type f)



build: frontend internal/gh db dist
	mkdir -p dist
	go build -o dist/onlyserve ./cmd/onlyserve

dev:
	saq -- bash -c 'make && ./dist/onlyserve --http localhost:8081'



dist: dist/schema.docs.graphql dist/static

dist/static: dist/static/script.js dist/static/styles.css

dist/static/script.js: $(SCRIPTS)
	./frontend/bundle.js ./frontend/scripts/index.js dist/static/script.js

dist/static/styles.css: frontend/styles.scss $(STYLES)
	sass ./frontend/styles.scss ./dist/static/styles.css

dist/schema.docs.graphql:
	wget -O dist/schema.docs.graphql https://gist.githubusercontent.com/diamondburned/6913a10f8c4ab97fbe341002b9d57840/raw/a610b7dd1fe7fb10e10680d132e5eb30ff6920b1/schema.docs.graphql

frontend: $(TEMPLS)
	cd frontend && templ generate
	touch frontend

internal/gh: internal/gh/queries.graphql.gen.go

internal/gh/queries.graphql.gen.go: internal/gh/queries.graphql genqlient.yaml dist/schema.docs.graphql
	genqlient

db: db/sqlitec/db.gen.go

db/sqlitec/db.gen.go: db/sqlc.json db/sqlitec/*.sql
	cd db && sqlc generate
