#!/bin/bash
rm -r ~/.qnscli
rm -r ~/.qnsd

qnsd init mynode --chain-id qns

qnscli config keyring-backend test
qnscli config chain-id qns
qnscli config output json
qnscli config indent true
qnscli config trust-node true

qnscli keys add user1
qnscli keys add user2
qnsd add-genesis-account $(qnscli keys show user1 -a) 1000token,100000000stake
qnsd add-genesis-account $(qnscli keys show user2 -a) 500token

qnsd gentx --name user1 --keyring-backend test

qnsd collect-gentxs 
