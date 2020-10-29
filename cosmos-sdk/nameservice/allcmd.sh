#!/bin/bash

nameservicecli query account $(nameservicecli keys show user1 -a) | jq ".value.coins[0]"
nameservicecli query account $(nameservicecli keys show user2 -a) | jq ".value.coins[0]"

# Buy your first name using your coins from the genesis file
nameservicecli tx nameservice buy-name youngqq.cn  5token --from user1 -y | jq ".txhash" |  xargs $(sleep 10) nameservicecli q tx

# Set the value for the name you just bought
nameservicecli tx nameservice set-name youngqq.cn 8.8.8.8 --from user1 -y | jq ".txhash" |  xargs $(sleep 10) nameservicecli q tx

# Try out a resolve query against the name you registered
nameservicecli query nameservice resolve-name youngqq.cn | jq ".value"
# > 8.8.8.8

# Try out a whois query against the name you just registered
nameservicecli query nameservice get-whois youngqq.cn | jq
# > {"value":"8.8.8.8","owner":"cosmos1l7k5tdt2qam0zecxrx78yuw447ga54dsmtpk2s","price":[{"denom":"nametoken","amount":"5"}]}

# Alice buys name from jack
nameservicecli tx nameservice buy-name youngqq.cn 10token --from user2 -y | jq ".txhash" |  xargs $(sleep 10) nameservicecli q tx

# Alice decides to delete the name she just bought from jack
nameservicecli tx nameservice delete-name youngqq.cn --from user2 -y | jq ".txhash" |  xargs $(sleep 10) nameservicecli q tx

# Try out a whois query against the name you just deleted
nameservicecli query nameservice get-whois youngqq.cn
# > {"value":"","owner":"","price":[{"denom":"nametoken","amount":"1"}]}