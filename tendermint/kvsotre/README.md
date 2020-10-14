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

提交一个交易

```bash
curl -s 'localhost:26657/broadcast_tx_commit?tx="tendermint=rocks"'
```

查询交易

```bash
curl -s 'localhost:26657/abci_query?data="tendermint"'
```


通过交易hash查询交易
https://docs.tendermint.com/master/rpc/#/Info/tx

```bash
curl -s "localhost:26657/tx?hash=0x9908D6B7D3213E16E722843048AF0C3E7CA7C3377E15C25FC349B321550F544A&prove=true"
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


