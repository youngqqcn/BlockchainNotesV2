# kvsotre

## 安装tendermint :
- tendermint core 环境

> tendermint安装: https://docs.tendermint.com/master/introduction/install.html

## 初始化

``` 
cd kvstore

tendermint init --home ./ 

```


### tendermint RPC 


> https://docs.tendermint.com/master/rpc/


```bash

# release  发币
curl -s 'localhost:26657/broadcast_tx_commit?tx="releasehello,yqq,1000"'

# transfer 转账
curl -s 'localhost:26657/broadcast_tx_commit?tx="yqq,tom,10"'

# query balance 查询账户余额
curl -s 'localhost:26657/abci_query?data="tom"'

```

https://docs.tendermint.com/master/rpc/#/Info/tx

```bash
curl -s localhost:26657/tx?hash=0x381E1C108E1CFF75CD7BA4CECE4ACE91CA46745C1E73DB47AD09E2C40CB97B61
```


获取区块
```bash
curl localhost:26657/block?height=12


curl localhost:26657/blockchain?minHeight=1&maxHeight=2
```


使用Websocket订阅event

```bash

wscat -c ws://127.0.0.1:26657/websocket

{"jsonrpc": "2.0","method": "subscribe","id": 0,"params": {"query": "tm.event='NewBlock'"}}

```



### V2

super user 
``` 
user6, 
private key: e1b0f79b20e4543e19dca67fa154746e09c63084ded8a4783861a7351dd69ec997eba1726b, 
address : 365EA5222D2F08A8A1EBF992B0628B1459527400
```





# 多节点

按照教程的方式改, 一直报错, 需要修改每个节点的  config.toml 

``` 
allow_duplicate_ip = true
addr_book_strict = false
```

