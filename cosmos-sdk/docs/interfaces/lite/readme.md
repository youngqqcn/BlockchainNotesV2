---
parent:
  order: false
---

# Light Client Overview

轻节点

**See the Cosmos SDK Light Client RPC documentation [here](https://cosmos.network/rpc/)**

## Introduction

A light client allows clients, such as mobile phones, to receive proofs of the state of the
blockchain from any full node. Light clients do not have to trust any full node, since they are able
to verify any proof they receive.

A light client can provide the same security as a full node with minimal requirements for
bandwidth, computing and storage resource. It can also provide modular(模块化的) functionality
according to users' configuration. These features allow developers to build secure, efficient,
and usable mobile apps, websites, and other applications without deploying or
maintaining any full blockchain nodes.

轻节点可以允许像手机这样的客户端, 接受全节点的区块链的状态证明. 轻节点不需要信任任何全节点, 因为他们能够验证收到的数据(merkle proof)

轻节点可以提供和全节点一样个的安全性, 而且不需要全节点那样的高的配置. 根据用户的配置可以提供模块化的功能. 它允许开发者构建安全,高效可用于移动app,网站的,以及其他不需要维护全节点的应用场景.


### What is a Light Client?

轻客户端

The Cosmos SDK Light Client (Gaia-lite) is split into two separate components. The first component is generic for
any Tendermint-based application. It handles the security and connectivity aspects of following the header
chain and verify proofs from full nodes against a locally trusted validator set. Furthermore, it exposes the same
API as any Tendermint Core node. The second component is specific for the Cosmos Hub (`gaiad`). It works as a query
endpoint and exposes the application specific functionality, which can be arbitrary. All queries against the
application state must go through the query endpoint. The advantage of the query endpoint is that it can verify
the proofs that the application returns.

Cosmos SDK 轻客户端分为两部分: 第一部分是基于Tendermint的应用程序, 它处理安全性和连接相关,验证信任的验证节点, 还提供了Tendermint相关的API. 第二部分, 是cosmos hub相关的.提供年一些查询接口,查询接口可以对返回数据进行验证.


### High-Level Architecture

An application developer that wants to build a third party client application for the Cosmos Hub (or any
other zone) should build it against its canonical API. That API is a combination of multiple parts.
All zones have to expose ICS0 (TendermintAPI). Beyond that any zone is free to choose any
combination of module APIs, depending on which modules the state machine uses. The Cosmos Hub will
initially support [ICS0](https://cosmos.network/rpc/#/ICS0) (TendermintAPI), [ICS1](https://cosmos.network/rpc/#/ICS1) (KeyAPI), [ICS20](https://cosmos.network/rpc/#/ICS20) (TokenAPI), [ICS21](https://cosmos.network/rpc/#/ICS21) (StakingAPI),
[ICS22](https://cosmos.network/rpc/#/ICS22) (GovernanceAPI) and [ICS23](https://cosmos.network/rpc/#/ICS23) (SlashingAPI).

![high-level](./pics/high-level.png)

All applications are expected to run only against Gaia-lite. Gaia-lite is the only piece of software
that offers stability guarantees around the zone API.

### Comparison

A full node of ABCI is different from a light client in the following ways:

|| Full Node | Gaia-lite | Description|
|-| ------------- | ----- | -------------- |
| Execute and verify transactions|Yes|No|A full node will execute and verify all transactions while Gaia-lite won't.|
| Verify and save blocks|Yes|No|A full node will verify and save all blocks while Gaia-lite won't.|
| Consensus participation|Yes|No|Only when a full node is a validator will it participate in consensus. Lite nodes never participate in consensus.|
| Bandwidth cost|High|Low|A full node will receive all blocks. If bandwidth is limited, it will fall behind the main network. What's more, if it happens to be a validator, it will slow down the consensus process. Light clients require little bandwidth, only when serving local requests.|
| Computing resources|High|Low|A full node will execute all transactions and verify all blocks, which requires considerable computing resources.|
| Storage resources|High|Low|A full node will save all blocks and ABCI states. Gaia-lite just saves validator sets and some checkpoints.|
| Power consumption|High|Low|Full nodes must be deployed on machines which have high performance and will be running all the time. Gaia-lite can be deployed on the same machines as users' applications, or on independent machines but with lower performance. Light clients can be shut down anytime when necessary. Gaia-lite consumes very little power, so even mobile devices can meet the power requirements.|
| Provide APIs|All cosmos APIs|Modular APIs|A full node supports all Cosmos APIs. Gaia-lite provides modular APIs according to users' configuration.|
| Secuity level| High|High|A full node will verify all transactions and blocks by itself. A light client can't do this, but it can query data from other full nodes and verify the data independently. Therefore, both full nodes and light clients don't need to trust any third nodes and can achieve high security.|

According to the above table, Gaia-lite can meet many users' functionality and security requirements, but require little bandwidth, computing, storage, and power.

## Achieving Security

### Trusted Validator Set

The base design philosophy of Gaia-lite follows two rules:

1. **Doesn't trust any blockchain nodes, including validator nodes and other full nodes**
2. **Only trusts the whole validator set**

- 不要信任任何区块链节点,验证节点和其他节点.
- 只信任整个验证节点集合

The original trusted validator set should be prepositioned(预设) into its trust store. Usually this
validator set comes from a genesis file. During runtime, if Gaia-lite detects a different validator set,
it will verify it and save the new validated validator set to the trust store.

![validator-set-change](./pics/validatorSetChange.png)

### Trust Propagation(信任传播)

From the above section, we come to know how to get a trusted validator set and how lcd keeps track of
validator set evolution(演化). The validator set is the foundation of trust, and the trust can propagate(传播) to
other blockchain data, such as blocks and transactions. The propagation architecture is shown as

follows:

![change-process](./pics/trustPropagate.png)

In general, with a trusted validator set, a light client can verify each block commit which contains all pre-commit
data and block header data. Then the block hash, data hash and appHash are trusted. Based on this
and merkle proof, all transactions data and ABCI states can be verified too.



通常, 有了可信任的验证节点集合, 轻客户端可以验证每个block commit. 因为区块中包含了pre-commit数据和区块头数据. 区块hash 和数据hash, appHash 是可以信任的. 基于这些和merkle proof, 所有的交易数据和ABCI状态都可以被验证.
