PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

EXED=nameserviced
EXECLI=nameservicecli
CHAINID=testchain

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=nameservice\
	-X github.com/cosmos/cosmos-sdk/version.ServerName=$(EXED)\
	-X github.com/cosmos/cosmos-sdk/version.ClientName=$(EXECLI)\
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

.PHONY: all
all: install

.PHONY: install


.PHONY: init

init:
	-rm -r ~/.$(EXED)
	-rm -r ~/.$(EXECLI)
	@$(EXED) init mynode --chain-id $(CHAINID)
	@$(EXECLI) config keyring-backend test
	@$(EXECLI) config chain-id $(CHAINID)
	@$(EXECLI) config output json
	@$(EXECLI) config indent true
	@$(EXECLI) config trust-node true
	@$(EXECLI) keys add user1
	@$(EXECLI) keys add user2
	@$(EXED) add-genesis-account $$($(EXECLI) keys show user1 -a) 1000token,100000000stake
	@$(EXED) add-genesis-account $$($(EXECLI) keys show user2 -a) 500token
	@$(EXED) gentx --name user1 --keyring-backend test
	@$(EXED) collect-gentxs 


start:
	@$(EXED) start  

start-rest:
	@$(EXECLI) rest-server

build:
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/$(EXED)
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/$(EXECLI)


install:  go.sum
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/$(EXED)
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/$(EXECLI)
		npm --prefix ./vue i --save 
		npm --prefix ./vue run build

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

# Uncomment when you have some tests
test:
	@go test -mod=readonly $(PACKAGES)
.PHONY: lint
# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify
