all: cleanall init build_mytokenapp build_client


init: cleanall
	@tendermint init --home ./bin

build_mytokenapp:
	@go build -o ./bin/mytokenapp main.go

build_client:
	@go build -o ./bin/client client.go

cleanbin:
	@rm -rf ./bin/mytoken*

cleanall:
	@rm -rf ./bin/*

start:
	@cd ./bin && ./mytokenapp


unitest:
	@cd mytokenapp && go test -v -count=1