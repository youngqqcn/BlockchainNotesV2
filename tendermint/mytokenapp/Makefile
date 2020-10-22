all: cleanall build_mytokenapp build_client


init:
	@cd ./bin && ./mytokenapp init --home . --trace

build_mytokenapp:
	@go build -o ./bin/mytokenapp main.go
	@cp ./bin/mytokenapp ./bin/n1/mytokenapp
	@cp ./bin/mytokenapp ./bin/n2/mytokenapp
	@cp ./bin/mytokenapp ./bin/n3/mytokenapp

reset_all:
	@rm -rf  bin/n1/account.db
	@rm -rf  bin/n2/account.db
	@rm -rf  bin/n3/account.db
	@cd ./bin/n1 && ./mytokenapp unsafe_reset_all && cd ..
	@cd ./bin/n2 && ./mytokenapp unsafe_reset_all && cd ..
	@cd ./bin/n3 && ./mytokenapp unsafe_reset_all && cd ..

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
