postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=jg.5s0o4ht7 -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16 dropdb simple_bank

gooseconfig:
	export GOOSE_DBSTRING="host=localhost port=5432 user=root dbname=simple_bank password=jg.5s0o4ht7 sslmode=disable"
	export GOOSE_DRIVER=postgres

.PHONY: postgres createdb dropdb gooseconfig
