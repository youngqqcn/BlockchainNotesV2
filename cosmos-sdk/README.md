# Cosmos-SDK 笔记


## Cosmos-SDK 的架构

> https://github.com/cosmos/cosmos-sdk/blob/master/docs/intro/sdk-app-architecture.md



区块链核心就是一个 Replicated Deterministic State Machine (可复制的确定性状态机)

### Tendermint

Tendermint 是一个与应用程序(区块链)无关的网络共识引擎, 基于拜占庭容错算法(Byzantine-Fault-Tolerant, BFT)并使用相同的交易顺序实现共识. 

Tendermint 框架中参与共识的节点成为 *Validators* (验证节点).验证者们负责选举下一个区块的提议者(proposer),并对提议者产生的新的区块进行验证和[投票, 如果投票数超过`2/3`的验证者数, 则区块是有效的,反之此次提议的区块无效,进行下一轮提议和投票, 这一轮提议和投票分为 `provote`  和  [`promcommit`](https://docs.tendermint.com/master/spec/consensus/consensus.html#precommit-step-height-h-round-r) 两个阶段. 当一个区块通过验证者验证之后, 那么就意味着此区块包含的交易都是有效的, 状态机就会进行更改, 同时更换下一个区块的提议者.



```
                ^  +-------------------------------+  ^
                |  |                               |  |   Built with Cosmos SDK
                |  |  State-machine = Application  |  |
                |  |                               |  v
                |  +-------------------------------+
                |  |                               |  ^
Blockchain node |  |           Consensus           |  |
                |  |                               |  |
                |  +-------------------------------+  |   Tendermint Core
                |  |                               |  |
                |  |           Networking          |  |
                |  |                               |  |
                v  +-------------------------------+  v
```


#### ABCI

ABCI (Application Blockchain  Interface)

Tendermint 只会处理交易字节(transaction bytes), 它不知道交易的具体内容是什么.所有的Tendermint节点处理的交易字节顺序都是确定的. Tendermint和应用程序通过ABCI进行数据交互, Tendermint通过ABCI将交易字节发给应用程序,只需知道返回的返回码是成功的还是失败的.

```
              +---------------------+
              |                     |
              |     Application     |
              |                     |
              +--------+---+--------+
                       ^   |
                       |   | ABCI
                       |   v
              +--------+---+--------+
              |                     |
              |                     |
              |     Tendermint      |
              |                     |
              |                     |
              +---------------------+
```

以下是几个重要的ABCI消息:

- `CheckTx`: 当Tendermint Core 接收到一笔交易, 这笔交易会被传递给*应用程序* 去检查交易的是否符合基本要求. `CheckTx`是用来保护全节点的交易池(mempool)抵抗垃圾交易攻击. 在CosmosSDK中有一个特殊的handler被称为[`AnteHandler`](https://github.com/cosmos/cosmos-sdk/blob/master/docs/basics/gas-fees.md#antehandler)用来执行一系列的交易检查操作, 例如检查交易费的有效性,交易签名的有效性. 如果检查通过, 这笔交易将被加入到交易池(mempool)然后通过p2p转发到其他节点. 注意:经过`CheckTx`,这些交易还没有被执行(尚未修改任何状态), 因为这些交易还没有被打包进区块中.

- `DeliverTx`: 当Tendermint Core 接收到一个有效的区块, 区块中的每个交易都通过`DeliverTx`*有序*地传递给应用程序执行. 在这个阶段状态机会发生改变. `AnteHandler` 与事务中的每个消息的实际处理程序一起再次执行.

- `BeginBlock/EndBlock` : 不管区块中是否包含交易, 在区块的开始和结束时这两个消息都会被执行. 这对于跟踪执行逻辑很有用. 但是需要谨慎，因为计算量大的循环可能拖慢区块链, 甚至无限循环会导致区块链停止。



## CosmosSDK 设计 - Cosmos SDK 主要组件

> https://github.com/cosmos/cosmos-sdk/blob/master/docs/intro/sdk-design.md

CosmosSDK是建立在Tendermint之上可以用来可开发安全的状态机的框架. Cosmos使用Golang实现ABCI. 拥有一个 `multistore` 进行数据持久化, `router`进行交易处理.


以下是一个简单的概述, 描述了一个基于CosmosSDK构建的应用程序是如何处理来自Tendermint通过`DeliverTx`发来的交易:

- 1. 解码(Decode)接受到的来自Tendermint共识引擎的`transactions`(记住Tendermint仅仅只处理交易字节 `[]byte` 不关心交易内容).
- 2. 提取(Extract) `transactions`中的 `messages`并且做一些基本的有效性检查.
- 3. 路由(Route)每个消息到相应的模块(module)进行消息处理.
- 4. 提交(Commit)状态的修改(state changes).


### `baseapp`


`baseapp` 是Cosmos SDK应用程序的一个样板实现.它带有一个用于处理和地层共识引擎的连接的ABCI实现. 通常, 一个Cosmos SDK应用程序会扩展 `baseapp`, 通过在`app.go`嵌入 `baseapp`.  这是一个来自Cosmos SDK教程中实际的例子:  https://github.com/cosmos/sdk-tutorials/blob/c6754a1e313eb1ed973c5c91dcc606f2fd288811/app.go#L72-L92

```go
type nameServiceApp struct {
	*bam.BaseApp           // 嵌入的 baseapp
	cdc *codec.Codec

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// Keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	distrKeeper    distr.Keeper
	supplyKeeper   supply.Keeper
	paramsKeeper   params.Keeper
	nsKeeper       nameservice.Keeper

	// Module Manager
	mm *module.Manager
}
```

`baseapp`的目标是在*存储(store)* 和 *可扩展的状态机(extensible state machine)* 提供一个安全的接口, 同时在ABCI保持不变的情况下尽可能地少定义状态机.

更多关于 `baseapp`, 可看[这里](https://github.com/cosmos/cosmos-sdk/blob/master/docs/core/baseapp.md)


### `Multistore`

Cosmos SDK 提供了一个名为 `multistore`组件用于状态的持久化(persisting state). `multistore`允许开发者声明任意数量的`KVStores`. 这些`KVStores`只能接收 `[]byte`类型作为值, 因此开发自定义的结构体在存储前需要使用一个[`codec`进行序列化](https://github.com/cosmos/cosmos-sdk/blob/master/docs/core/encoding.md). 

抽象`multistore`用于将不同模块的状态分开, 不同的模块只需要维护自己的状态 .关于`multisotre`更多的详情, [点击这里](https://github.com/cosmos/cosmos-sdk/blob/master/docs/core/store.md#multistore)



### Modules

Cosmos SDK 强大之处在于它的模块化. 应用程序使用一些SDK的*可互相操作*(interoperate)的模块集合进行构建. 每个模块定义了一个状态子集并且包含了模块自己的`message/transaction`处理器(processor), 而SDK负责将每个消息(message)路由到它们各自的模块.


以下是继续Cosmos SDK构建的应用程序的全节点收到一个有效区块后处理一笔交易的大致流程:

```
                                      +
                                      |
                                      |  Transaction relayed from the full-node's Tendermint engine
                                      |  to the node's application via DeliverTx
                                      |  
                                      |
                                      |
                +---------------------v--------------------------+
                |                 APPLICATION                    |
                |                                                |
                |     Using baseapp's methods: Decode the Tx,    |
                |     extract and route the message(s)           |
                |                                                |
                +---------------------+--------------------------+
                                      |
                                      |
                                      |
                                      +---------------------------+
                                                                  |
                                                                  |
                                                                  |
                                                                  |  Message routed to the correct
                                                                  |  module to be processed
                                                                  |
                                                                  |
+----------------+  +---------------+  +----------------+  +------v----------+
|                |  |               |  |                |  |                 |
|  AUTH MODULE   |  |  BANK MODULE  |  | STAKING MODULE |  |   GOV MODULE    |
|                |  |               |  |                |  |                 |
|                |  |               |  |                |  | Handles message,|
|                |  |               |  |                |  | Updates state   |
|                |  |               |  |                |  |                 |
+----------------+  +---------------+  +----------------+  +------+----------+
                                                                  |
                                                                  |
                                                                  |
                                                                  |
                                       +--------------------------+
                                       |
                                       | Return result to Tendermint
                                       | (0=Ok, 1=Err)
                                       v
```

每个module都可以被看作是一个小的状态机. 开发者需要定义模块处理的状态子集(subset of state), 还要自定义修改状态的消息(message)类型(注意: `message`是通过`baseapp`从`transactions`提取出来的). 通常, 每个模块在`multistore`定义它自己的`KVStore`来存储这些模块定义的状态子集. 大多数开发者在构建自己的模块时需要访问其他第三方的模块. 鉴于Cosmos-SDK是一个开放框架,某些模块可能是恶意的,这意味着需要安全原则来规范模块间的交互. 这个原则基于[object-capabilities](https://github.com/cosmos/cosmos-sdk/blob/master/docs/core/ocap.md). 实际上, 并不是说着要为其他模块保存一个访问控制列表, 而是每个模块实现了一个名为`keepers`的特殊对象(speicial objects), `keepers`可以被传递给其他模块以赋予其一组预定义(pre-defined)的功能.


SDK 的模块定义在 `x/` 目录下, 以下是几个核心的模块:

- `x/auth` : 用于管理账户和签名
- `x/bank` : 用于激活代币和代币转账
- `x/staking` + `x/slashing`: 用于构建权益证明`POS`(Proof-Of-Stake)区块链

 除了可以使用`x/`目录中已经存在的模块之外 , SDK还允许你构建自己的模块用于自己的应用程序中.



## Cosmos SDK 核心组件剖析

### Node Client

守护进程或者说全节点, 是基于SDK的区块链的核心进程. 全节点进程用于初始化状态机,连接其他全节点并在一个接收到一个新的区块时更新状态机.

```
                ^  +-------------------------------+  ^
                |  |                               |  |
                |  |  State-machine = Application  |  |
                |  |                               |  |   Built with Cosmos SDK
                |  |            ^      +           |  |
                |  +----------- | ABCI | ----------+  v
                |  |            +      v           |  ^
                |  |                               |  |
Blockchain Node |  |           Consensus           |  |
                |  |                               |  |
                |  +-------------------------------+  |   Tendermint Core
                |  |                               |  |
                |  |           Networking          |  |
                |  |                               |  |
                v  +-------------------------------+  v
```

区块链全节点是可执行文件, 通常以 `d`(daemon)结尾. 这个可执行文件通过编译位于`./cmd/appd/`目录下的`main.go`得到. 这个操作一般写在 Makefile里面.

一旦主的可执行文件被构建, 就可以通过 `start` 命令进行启动.这个命令的主要作用是做以下3件事情:

- 1. 创建一个在`app.go`中的定义的状态机实例.
- 2. 使用最新的已知的状态初始化状态机, 从存储在`~/.appd/data`目录提取`db`. 此时, 状态机的高度是 `appBlockHeight`.
- 3. 创建并启动一个新的Tendermint实例. 除其他事项外, 节点会执行与对端节点(peers)的握手, 从其他节点获取最新的区块高度 `blockHeight`, 如果最新区块大于本地的最高区块高度`appBlockHeight`, 那么,节点会重放区块(replay blocks)以同步到最新区块高度. 如果`appBlockHeight`是`0`, 那么节点会从创始区块开始启动, 并且Tendermint会通过ABCI发送一个`InitChain`消息给应用程序`app`, 它会触发 `InitChainer`.



### Core Application File

一般状态机的核心定义在一个名为`app.go`的文件中.它主要包含了应用程序的类型定义以及创建,初始化的一些函数.

#### 应用程序的类型定义


- **`baseapp`的引用(指针)**: 在`app.go`中定义的应用是`baseapp`的一个扩展(子类). 当一个交易被Tendermint中继到应用程序时, `app`会使用`baseapp`的方法路由到对应的模块. `baseapp`实现了应用程序的大部分核心逻辑, 包括所有的ABCI方法和路由逻辑(routing logic).

- **一组key stores**: `store`包含了完整的状态, 在Cosmos SDK 是以`multistore`( a store of stores ) 实现的. 每个模块使用一个或多个`multisotre`中的`stores`保存各自模块中的状态. 这些`stores`可以通过声明在`app`类型中的key进行访问. 这些 key和 `keepers` 是Cosmos SDK的 `object-capabilities`(对象能力)模型的核心.

- **一组模块的`keeper`**:  每个模块定义了一个名为`keeper`抽象, 用于处理模块的store(s)的读和写. 一个模块的`keeper`的方法可以被其他模块调用(如果被授权), 之所以要把他们定义在应用程序类型中, 并且导出为interface给其他模块, 是为了其他模块能够访问这些已授权方法.

- **一个`appCodec`的引用(指针)** : `appCodec`是用于对数据结构进行序列化和反序列化以便进行存储, 因为stores只能使用`[]byte`进行存储.默认的codec是 Protocol Buffers.

- **一个`legacyAmino` 的引用(指针) codec**: SDK的某些部分尚未迁移至上述`appCodec`, 依然使用使用Amino进行硬编码. 其他部分也显示的使用Amino以保持向后兼容.基于这些原因, 应用程序依然持有一个传统的Amino解码器(codec). 但是请注意 Amino codec将在将来的发行版中移除.


- **一个moudle manager的引用和一个 basic moudle manager**:  模块管理器(moudle manager)是一个包含一组应用程序模块的对象. 它用于对模块进行一些操作, 如 注册路由(registering `routes`), gRPC 查询服务和Tendermint基本查询路由 或者设置不同模块之间的函数执行顺序,如`InitChainer`, `BeginBlocker`和`EndBlocker`.


如下`simapp` [是SDK自己的用于测试的demo](https://github.com/cosmos/cosmos-sdk/blob/d9175200920e96bfa4182b5c8bc46d91b17a28a1/simapp/app.go#L140-L179)

```go

// SimApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type SimApp struct {
	*baseapp.BaseApp
	cdc               *codec.LegacyAmino
	appCodec          codec.Marshaler
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper   evidencekeeper.Keeper
	TransferKeeper   ibctransferkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedIBCMockKeeper  capabilitykeeper.ScopedKeeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}
```


#### 构造者函数(Constructor Function)

这个函数用于构造一个新的如上节所定义的`Application`实例, 必须符合 `AppCreator`函数签名, 以便在 `start`命令中使用.


```go

// AppCreator is a function that allows us to lazily initialize an
// application using various configurations.
AppCreator func(log.Logger, dbm.DB, io.Writer, AppOptions) Application

```

以下是构造者函数的主要执行的操作:


- 实例化一个新的`codec`并且使用`basic manager`初始化应用程序每个模块的`codec`
- 实例化一个新的应用程序对象, 其包含`baseapp`引用的和一个`codec`和其他的相应的store keys.
- 使用应用程序的每个模块的`NewKeeper`函数实例化所有的在应用程序类型中定义的`keeper`. 注意: `keepers`必须以正确的顺序进行实例化, 因为一个模块的`NewKeeper`可能需要引用其他模块的`keeper`.
- 使用应用程序的每个模块的`AppModule`对象实例化应用程序的 `module manager`.
- 使用 `module manager`初始化应用程序的 `routes`, `gRPC query serives` 和 `legacy query routes`. 当一个交易经由Tendermint通过ABCI发送给应用程序时, 使用`routes`将交易路由到相应模块的`handler`. 同样地, 当一个应用程序受到一个gRPC请求时, 会路由到相应的`gRPC query service`. SDK依然支持传统的Tendermint查询, 这些查询都是通过`legacy query routes`进行路由.

- 使用`module manager`注册应用程序模块的的不变量(application's modules' invariants). 例如, token总发行量. 它会在每个区块中进行检查. 检查不变量的操作是通过一个名为 `InvariantRegistry`的模块进行. 不变量的值应该等于模块中定义的值.如果该值与预测的值不同,那么在不变量注册表中定义的操作将本触发(一般表现为区块链挂掉). 早日发现问题, 这对于确保不会发现任何严重错误并产生难以修复的长期影响非常有用.
- 使用`module manager` 设置应用程序的每个模块的 `InitGenesis`, `BeginBlocker` 和`EndBlocker`的执行顺序. 注意并不是所有的模块都实现了这些函数.
- 设置应用程序其他的参数:
  - `InitChianer`: 被用于第一次启动时初始化应用程序
  - `BeginBlocker`, `EndBlocker`: 在开始区块和结束区块是被调用
  - `anteHandler`: 被用于处理验证手续费和签名有效性

- 挂载 stores
- 返回应用程序实例

注意:这个函数仅创建了一个应用实例, 而实际的状态将从 `~/.appd/data`目录中加载, 如果是第一次启动, 则从创始文件(`genesis file`)生成.


可以看`simapp`示例代码: https://github.com/cosmos/cosmos-sdk/blob/d9175200920e96bfa4182b5c8bc46d91b17a28a1/simapp/app.go#L190-L427


#### InitChainer


