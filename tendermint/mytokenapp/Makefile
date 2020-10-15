all: cleanall init build_mytokenapp build_client


init: cleanall
	@tendermint init --home ./bin

build_mytokenapp:
	@go build -o ./bin/mytokenapp main.go

build_client: clean_client
	@go build -o ./bin/client client.go

clean_client:
	-rm -f ./bin/client

clean_mytokenapp:
	-rm -f ./bin/mytokenapp

cleanall:
	-rm -rf ./bin/*

start:
	@cd ./bin && ./mytokenapp

unitest:
	@cd mytokenapp && go test -v -count=1

initwallet: build_client
	@./bin/client  initwallet --filepath ./bin/wallet.dat -lsuperuser -lyqq -ltom -ljack -lalice -lbob

releasetoken:
	@cd ./bin && ./client release -tyqq -v9996

transfertoken:
	@cd ./bin && ./client transfer -fyqq -ttom -v11

querybalance:
	@cd ./bin &&  ./client querybalance -w wallet.dat -lyqq