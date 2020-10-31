<!--
order: 3
-->

# Context

The `context` is a data structure intended to be passed from function to function that carries information about the current state of the application. It holds a cached copy of the entire state as well as useful objects and information like `gasMeter`, `block height`, `consensus parameters` and more. {synopsis}

`context`是一个用于函数间传递应用程序当前状态的信息的数据结构. 它保存了完整状态的一个缓存的备份和一些有用的对象信息, 以及类似 `gasMeter`, `blocl height` `consensus parameters`等信息.

## Pre-requisites Readings

- [Anatomy of an SDK Application](../basics/app-anatomy.md) {prereq}
- [Lifecycle of a Transaction](../basics/tx-lifecycle.md) {prereq}

## Context Definition

The SDK `Context` is a custom data structure that contains Go's stdlib [`context`](https://golang.org/pkg/context) as its base, and has many additional types within its definition that are specific to the Cosmos SDK. The `Context` is integral to transaction processing in that it allows modules to easily access their respective [store](./store.md#base-layer-kvstores) in the [`multistore`](./store.md#multistore) and retrieve transactional context such as the block header and gas meter.


SDK `Context`是自定义的数据结构, 允许模块方便地访问store, 并且获取交易的context,如 block header 和 gas meter

```go
type Context struct {
  ctx           context.Context   // go  stdlib  context
  ms            MultiStore  // 应用程序的 KVStore 和 TransientStore
  header        tmproto.Header  // 保存了区块连很重要的状态信息, 例如 区块高度,  区块提议者
  chainID       string 
  txBytes       []byte // 交易字节
  logger        log.Logger
  voteInfo      []abci.VoteInfo  // 包含了验证节点的名称
  gasMeter      GasMeter
  blockGasMeter GasMeter
  checkTx       bool //一笔交易是否检查
  minGasPrice   DecCoins   //(本地的)最小gas价格
  consParams    *abci.ConsensusParams  // 共识参数
  eventManager  *EventManager  //事件管理器,
}
```

- **Context:** The base type is a Go [Context](https://golang.org/pkg/context), which is explained further in the [Go Context Package](#go-context-package) section below. 
- **Multistore:** Every application's `BaseApp` contains a [`CommitMultiStore`](./store.md#multistore) which is provided when a `Context` is created. Calling the `KVStore()` and `TransientStore()` methods allows modules to fetch their respective [`KVStore`](./store.md#base-layer-kvstores) using their unique `StoreKey`.
- **ABCI Header:** The [header](https://tendermint.com/docs/spec/abci/abci.html#header) is an ABCI type. It carries important information about the state of the blockchain, such as block height and proposer of the current block.
- **Chain ID:** The unique identification number of the blockchain a block pertains to.
- **Transaction Bytes:** The `[]byte` representation of a transaction being processed using the context. Every transaction is processed by various parts of the SDK and consensus engine (e.g. Tendermint) throughout its [lifecycle](../basics/tx-lifecycle.md), some of which to not have any understanding of transaction types. Thus, transactions are marshaled into the generic `[]byte` type using some kind of [encoding format](./encoding.md) such as [Amino](./encoding.md).
- **Logger:** A `logger` from the Tendermint libraries. Learn more about logs [here](https://tendermint.com/docs/tendermint-core/how-to-read-logs.html#how-to-read-logs). Modules call this method to create their own unique module-specific logger.
- **VoteInfo:** A list of the ABCI type [`VoteInfo`](https://tendermint.com/docs/spec/abci/abci.html#voteinfo), which includes the name of a validator and a boolean indicating whether they have signed the block.
- **Gas Meters:** Specifically, a [`gasMeter`](../basics/gas-fees.md#main-gas-meter) for the transaction currently being processed using the context and a [`blockGasMeter`](../basics/gas-fees.md#block-gas-meter) for the entire block it belongs to. Users specify how much in fees they wish to pay for the execution of their transaction; these gas meters keep track of how much [gas](../basics/gas-fees.md) has been used in the transaction or block so far. If the gas meter runs out, execution halts.
- **CheckTx Mode:** A boolean value indicating whether a transaction should be processed in `CheckTx` or `DeliverTx` mode.
- **Min Gas Price:** The minimum [gas](../basics/gas-fees.md) price a node is willing to take in order to include a transaction in its block. This price is a local value configured by each node individually, and should therefore **not be used in any functions used in sequences leading to state-transitions**. 
- **Consensus Params:** The ABCI type [Consensus Parameters](https://tendermint.com/docs/spec/abci/apps.html#consensus-parameters), which specify certain limits for the blockchain, such as maximum gas for a block.
- **Event Manager:** The event manager allows any caller with access to a `Context` to emit [`Events`](./events.md). Modules may define module specific
`Events` by defining various `Types` and `Attributes` or use the common definitions found in `types/`. Clients can subscribe or query for these `Events`. These `Events` are collected throughout `DeliverTx`, `BeginBlock`, and `EndBlock` and are returned to Tendermint for indexing. For example:

```go
ctx.EventManager().EmitEvent(sdk.NewEvent(
    sdk.EventTypeMessage,
    sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory)),
)
```

## Go Context Package

A basic `Context` is defined in the [Golang Context Package](https://golang.org/pkg/context). A `Context`
is an immutable data structure that carries request-scoped data across APIs and processes. Contexts
are also designed to enable concurrency and to be used in goroutines.

`Context`是一个不可变的数据结构, 携带了一些接口请求和处理的所需的数据. `Context`的设计支持并发,可以用于`goroutine`s

Contexts are intended to be **immutable**; they should never be edited. Instead, the convention is
to create a child context from its parent using a `With` function. For example:

`Context`设计的原则就是**不可变**; 所以, 不要修改context. 可以创建子context:


``` go
childCtx = parentCtx.WithBlockHeader(header)
```

The [Golang Context Package](https://golang.org/pkg/context) documentation instructs developers to
explicitly pass a context `ctx` as the first argument of a process.

## Cache Wrapping

The `Context` contains a `MultiStore`, which allows for cache-wrapping functionality: a `CacheMultiStore`
where each `KVStore` is is wrapped with an ephemeral(短暂的) cache. Processes are free to write changes to
the `CacheMultiStore`, then write the changes back to the original state or disregard them if something
goes wrong. The pattern of usage for a Context is as follows:

`Context`包含了一个 `MultiStore`, 它支持 cache-wrapping(可以理解为副本功能): `CacheMultiStore`中每个`KVStore`都是一个临时缓存(副本). 可以自由对`CacheMultiStore`进行修改, 如果出错, 一切修改将丢弃; 如果一切正常,则将修改写回到原始的状态机中.


1. A process receives a Context `ctx` from its parent process, which provides information needed to
   perform the process.
2. The `ctx.ms` is **cache wrapped**, i.e. a cached copy of the [multistore](./store.md#multistore) is made so that the process can make changes to the state as it executes, without changing the original`ctx.ms`. This is useful to protect the underlying multistore in case the changes need to be reverted at some point in the execution. 
3. The process may read and write from `ctx` as it is executing. It may call a subprocess and pass
`ctx` to it as needed.
4. When a subprocess returns, it checks if the result is a success or failure. If a failure, nothing
needs to be done - the cache wrapped `ctx` is simply discarded. If successful, the changes made to
the cache-wrapped `MultiStore` can be committed to the original `ctx.ms` via `Write()`.



For example, here is a snippet from the [`runTx`](./baseapp.md#runtx-and-runmsgs) function in
[`baseapp`](./baseapp.md):

```go

// 1. 使用cacheTxContext将context和multisotr 拷贝一个副本
runMsgCtx, msCache := app.cacheTxContext(ctx, txBytes)
// 2. 使用副本执行消息
result = app.runMsgs(runMsgCtx, msgs, mode)
result.GasWanted = gasWanted


if mode != runTxModeDeliver {
  // 3.如果是 checkTxMode , 则直接返回即可
  return result
}

// 4.如果是deliverTxMode , 则需要修改状态, 即需要将副本中的修改, 写回到原来状态机中
if result.IsOK() { // 只有成功了才进行状态更改
  msCache.Write()
}

// 否则不会修改状态
```

Here is the process:

1. Prior to calling `runMsgs` on the message(s) in the transaction, it uses `app.cacheTxContext()`
to cache-wrap the context and multistore.
2. The cache-wrapped context, `runMsgCtx`, is used in `runMsgs` to return a result.
3. If the process is running in [`checkTxMode`](./baseapp.md#checktx), there is no need to write the
changes - the result is returned immediately.
4. If the process is running in [`deliverTxMode`](./baseapp.md#delivertx) and the result indicates
a successful run over all the messages, the cached multistore is written back to the original.

## Next {hide}

Learn about the [node client](./node.md) {hide}
