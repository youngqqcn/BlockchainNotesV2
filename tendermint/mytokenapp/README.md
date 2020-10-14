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
curl localhost:26657/blockchain?minHeight=1&maxHeight=2
```


使用Websocket订阅event

```bash

wscat -c ws://127.0.0.1:26657/websocket

{"jsonrpc": "2.0","method": "subscribe","id": 0,"params": {"query": "tm.event='NewBlock'"}}

```


