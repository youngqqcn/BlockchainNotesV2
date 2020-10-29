
# 获取账户信息
curl -s http://localhost:1317/auth/accounts/$(nameservicecli keys show user1 -a)

curl -s http://localhost:1317/auth/accounts/$(nameservicecli keys show user2 -a)

curl -X POST -s http://localhost:1317/nameservice/whois --data-binary '{"base_req":{"from":"'$(nameservicecli keys show user1 -a)'","chain_id":"testchain"},"name":"youngqq.cn","bid":"150token","buyer":"'$(nameservicecli keys show user1 -a)'"}' > unsignedTx.json

# 要注意 accountNumber 和 sequence
nameservicecli tx sign unsignedTx.json --from user1 --offline --chain-id testchain --sequence 1 --account-number 2 > signedTx.json

nameservicecli tx broadcast signedTx.json


curl -X PUT -s http://localhost:1317/nameservice/whois --data-binary '{"base_req":{"from":"'$(nameservicecli keys show user1 -a)'","chain_id":"testchain"},"name":"youngqq.cn","value":"8.8.4.4","owner":"'$(nameservicecli keys show user1 -a)'"}' > unsignedTx.json

# 要注意 accountNumber 和 sequence
nameservicecli tx sign unsignedTx.json --from user1 --offline --chain-id testchain --sequence 2 --account-number 0 > signedTx.json
nameservicecli tx broadcast signedTx.json

curl -s http://localhost:1317/nameservice/whois/resolve-name/youngqq.cn

curl -s http://localhost:1317/nameservice/whois/youngqq.cn


# 高价从 user1 手中买过来
curl -X POST -s http://localhost:1317/nameservice/whois --data-binary '{"base_req":{"from":"'$(nameservicecli keys show user2 -a)'","chain_id":"testchain"},"name":"youngqq.cn","bid":"200token","buyer":"'$(nameservicecli keys show user2 -a)'"}' > unsignedTx.json

nameservicecli tx sign unsignedTx.json --from user2 --offline --chain-id testchain --sequence 1 --account-number 1 > signedTx.json
nameservicecli tx broadcast signedTx.json


# 删除
curl -XDELETE -s http://localhost:1317/nameservice/whois --data-binary '{"base_req":{"from":"'$(nameservicecli keys show user2 -a)'","chain_id":"testchain"},"name":"youngqq.cn","owner":"'$(nameservicecli keys show user2 -a)'"}' > unsignedTx.json

curl -s http://localhost:1317/auth/accounts/$(nameservicecli keys show user2 -a)
# 注意 sequence 和 accountNumber
nameservicecli tx sign unsignedTx.json --from user2 --offline --chain-id testchain --sequence 2 --account-number 1 > signedTx.json
nameservicecli tx broadcast signedTx.json

curl -s http://localhost:1317/nameservice/whois/youngqq.cn