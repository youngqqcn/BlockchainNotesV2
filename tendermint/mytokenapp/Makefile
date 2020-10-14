all: cleanall init build


init: cleanall
	@tendermint init --home ./bin

build: cleanbin
	@go build -o ./bin/mytokenapp

cleanbin:
	@rm -rf ./bin/mytoken*

cleanall:
	@rm -rf ./bin/*

start:
	@cd ./bin && ./mytokenapp
