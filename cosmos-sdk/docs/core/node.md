<!--
order: 4
-->

# Node Client (Daemon)

The main endpoint of an SDK application is the daemon client, otherwise known as the full-node client. The full-node runs the state-machine, starting from a genesis file. It connects to peers running the same client in order to receive and relay transactions, block proposals and signatures. The full-node is constituted of the application, defined with the Cosmos SDK, and of a consensus engine connected to the application via the ABCI. {synopsis}

SDK应用程序的主要形式是守护进程, 也就是常说的全节点程序. 全节点运行着状态机, 从创世文件启动. 它会连接到对等节点, 以便接收和中继交易, 区块提议和签名.全节点由Cosmos SDK应用程序, 共识引擎通过ABCI与应用程序进行连接.

## Pre-requisite Readings

- [Anatomy of an SDK application](../basics/app-anatomy.md) {prereq}

## `main` function

The full-node client of any SDK application is built by running a `main` function. The client is generally named by appending the `-d` suffix to the application name (e.g. `appd` for an application named `app`), and the `main` function is defined in a `./cmd/appd/main.go` file. Running this function creates an executable `.appd` that comes with a set of commands. For an app named `app`, the main command is [`appd start`](#start-command), which starts the full-node. 

In general, developers will implement the `main.go` function with the following structure:

- First, a [`codec`](./encoding.md) is instanciated for the application.
- Then, the `config` is retrieved and config parameters are set. This mainly involves(涉及) setting the bech32 prefixes for [addresses and pubkeys](../basics/accounts.md#addresses-and-pubkeys).
	+++ https://github.com/cosmos/cosmos-sdk/blob/7d7821b9af132b0f6131640195326aa02b6751db/types/config.go#L10-L21
- Using [cobra](https://github.com/spf13/cobra), the root command of the full-node client is created. After that, all the custom commands of the application are added using the `AddCommand()` method of `rootCmd`. 
- Add default server commands to `rootCmd` using the `server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)` method. These commands are separated from the ones added above since they are standard and defined at SDK level. They should be shared by all SDK-based applications. They include the most important command: the [`start` command](#start-command).
- Prepare and execute the `executor`.  
	+++ https://github.com/tendermint/tendermint/blob/bc572217c07b90ad9cee851f193aaa8e9557cbc7/libs/cli/setup.go#L75-L78

See an example of `main` function from the [`gaia`](https://github.com/cosmos/gaia) application:

+++ https://github.com/cosmos/gaia/blob/f41a660cdd5bea173139965ade55bd25d1ee3429/cmd/gaiad/main.go

## `start` command

The `start` command is defined in the `/server` folder of the Cosmos SDK. It is added to the root command of the full-node client in the [`main` function](#main-function) and called by the end-user to start their node:

```go
// For an example app named "app", the following command starts the full-node

appd start
```

As a reminder, the full-node is composed of three conceptual layers: the networking layer, the consensus layer and the application layer. The first two are generally bundled together in an entity called the consensus engine (Tendermint Core by default), while the third is the state-machine defined with the help of the Cosmos SDK. Currently, the Cosmos SDK uses Tendermint as the default consensus engine, meaning the start command is implemented to boot up a Tendermint node. 


全节点由三个概念层: 网络层, 共识层, 应用层.   Tendermint负责了网络层和共识层.  应用层是一个使用Cosmos SDK构建的状态机.  


The flow of the `start` command is pretty straightforward(简单粗暴). First, it retrieves the `config` from the `context` in order to open the `db` (a [`leveldb`](https://github.com/syndtr/goleveldb) instance by default). This `db` contains the latest known state of the application (empty if the application is started from the first time. 

`start`启动命令的过程: 第一, 从context中获取`config`以便打开`db`. `db`包含了应用程序最新的状态(如果是第一次启动, 则是空的).

With the `db`, the `start` command creates a new instance of the application using an `appCreator` function:

有了`db`就可以创建应用程序实例了, 使用 `appCreator`

+++ https://github.com/cosmos/cosmos-sdk/blob/7d7821b9af132b0f6131640195326aa02b6751db/server/start.go#L144

Note that an `appCreator` is a function that fulfills the `AppCreator` signature. In practice, the [constructor the application](../basics/app-anatomy.md#constructor-function) is passed as the `appCreator`.

+++ https://github.com/cosmos/cosmos-sdk/blob/7d7821b9af132b0f6131640195326aa02b6751db/server/constructors.go#L17-L25

```go
type (
	// AppCreator is a function that allows us to lazily initialize an
	// application using various configurations.
	// 创建应用程序实例
	AppCreator func(log.Logger, dbm.DB, io.Writer) abci.Application

	// AppExporter is a function that dumps all app state to
	// JSON-serializable structure and returns the current validator set.
	AppExporter func(log.Logger, dbm.DB, io.Writer, int64, bool, []string) (json.RawMessage, []tmtypes.GenesisValidator, error)
)

```


Then, the instance of `app` is used to instanciate a new Tendermint node:
使用 `app`实例实例化一个 Tendermint 节点
+++ https://github.com/cosmos/cosmos-sdk/blob/7d7821b9af132b0f6131640195326aa02b6751db/server/start.go#L153-L163


```go
// create & start tendermint node
tmNode, err := node.NewNode(
	cfg,
	pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
	nodeKey,
	proxy.NewLocalClientCreator(app),  // 使用app实例
	node.DefaultGenesisDocProviderFunc(cfg),
	node.DefaultDBProvider,
	node.DefaultMetricsProvider(cfg.Instrumentation),
	ctx.Logger.With("module", "node"),
)
```

```go
// app应用程序需要实现的接口, Tendermint 会通过ABCI与应用程序通信(非go语言), 如果是使用go实现的, 则直接是通过函数调用进行通信
// Application is an interface that enables any finite, deterministic state machine
// to be driven by a blockchain-based replication engine via the ABCI.
// All methods take a RequestXxx argument and return a ResponseXxx argument,
// except CheckTx/DeliverTx, which take `tx []byte`, and `Commit`, which takes nothing.
type Application interface {
	// Info/Query Connection
	Info(RequestInfo) ResponseInfo                // Return application info
	SetOption(RequestSetOption) ResponseSetOption // Set application option
	Query(RequestQuery) ResponseQuery             // Query for state

	// Mempool Connection
	CheckTx(RequestCheckTx) ResponseCheckTx // Validate a tx for the mempool

	// Consensus Connection
	InitChain(RequestInitChain) ResponseInitChain    // Initialize blockchain w validators/other info from TendermintCore
	BeginBlock(RequestBeginBlock) ResponseBeginBlock // Signals the beginning of a block
	DeliverTx(RequestDeliverTx) ResponseDeliverTx    // Deliver a tx for full processing
	EndBlock(RequestEndBlock) ResponseEndBlock       // Signals the end of a block, returns changes to the validator set
	Commit() ResponseCommit                          // Commit the state and return the application Merkle root hash
}
```

The Tendermint node can be created with `app` because the latter satisfies the [`abci.Application` interface](https://github.com/tendermint/tendermint/blob/bc572217c07b90ad9cee851f193aaa8e9557cbc7/abci/types/application.go#L11-L26) (given that `app` extends [`baseapp`](./baseapp.md)). As part of the `NewNode` method, Tendermint makes sure that the height of the application (i.e. number of blocks since genesis) is equal to the height of the Tendermint node. The difference between these two heights should always be negative or null. If it is strictly negative, `NewNode` will replay blocks until the height of the application reaches the height of the Tendermint node. Finally, if the height of the application is `0`, the Tendermint node will call [`InitChain`](./baseapp.md#initchain) on the application to initialize the state from the genesis file. 

Once the Tendermint node is instanciated and in sync with the application, the node can be started:

```go
// 启动Tendermint节点
if err := tmNode.Start(); err != nil { 
	return nil, err
}
```

Upon starting, the node will bootstrap(启动) its RPC and P2P server and start dialing peers(P2P节点发现相关). During handshake with its peers, if the node realizes they are ahead, it will query all the blocks sequentially in order to catch up. Then, it will wait for new block proposals and block signatures from validators in order to make progress. 

## Next {hide}

Learn about the [store](./store.md) {hide}