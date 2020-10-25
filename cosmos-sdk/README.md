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

`InitChainer`是从创世文件(即创世账号(genesis accounts)的余额)初始化应用程序状态的函数. 当节点从高度0开始启动时(`appBlockHeight == 0`), 应用程序收到来自Tendermint的`InitChain`消息时会调用此函数. 应用程序必须在它的构造函数(constructor)中通过`SetInitChainer`方法设置 `InitChainer`.

通常, `InitChainer`主要是由每个应用程序模块的`InitGenesis`函数组成. 这可以通过调用`module manager`的`InitGenesis`函数来完成, `module manager`会依次调用每个模块包含的`InitGenesis`函数.注意模块的`InitGenesis`函数的被调用的顺序必须在`module manager`的`SetOrderInitGenesis`中进行设置. 这是在应用程序的构造函数中完成的, `SetOrderInitGenesis`必须在`SetInitChainer`之前被调用.

示例代码: https://github.com/cosmos/cosmos-sdk/blob/d9175200920e96bfa4182b5c8bc46d91b17a28a1/simapp/app.go#L450-L455




#### BeginBlocker and EndBlocker

该SDK为开发人员提供了在其应用程序中实现代码自动执行的可能性。这是通过两个称为`BeginBlocker`和`EndBlocker`的函数实现的。当应用程序分别从Tendermint引擎接收`BeginBlock`和`EndBlock`消息时，将调用它们，它们发生在每个块的开始和结尾。应用程序必须通过`SetBeginBlocker`和`SetEndBlocker`方法在其构造函数中设置`BeginBlocker`和`EndBlocker`.

通常，`BeginBlocker`和`EndBlocker`函数主要由每个应用程序模块的`BeginBlock`和`EndBlock`函数组成。这是通过调用模块管理器(module manager)的`BeginBlock`和`EndBlock`函数完成的，而后者又将调用其包含的每个模块的`BeginBLock`和`EndBlock`函数。请注意，必须在模块管理器中分别使用`SetOrderBeginBlock`和`SetOrderEndBlock`方法设置必须调用模块的`BegingBlock`和`EndBlock`函数的顺序。这是通过应用程序构造函数中的模块管理器完成的，必须在`SetBeginBlocker`和`SetEndBlocker`函数之前调用`SetOrderBeginBlock`和`SetOrderEndBlock`方法.

附带说明，请记住特定于应用程序的区块链是确定性的，这一点很重要。开发人员必须注意不要在`BeginBlocker`或`EndBlocker`中引入非确定性，并且还必须注意不要使它们在计算上过于昂贵，因为`gas`机制不会限制`BeginBlocker`和`EndBlocker`执行的成本。

查看simapp中的BeginBlocker和EndBlocker函数的示例:
https://github.com/cosmos/cosmos-sdk/blob/d9175200920e96bfa4182b5c8bc46d91b17a28a1/simapp/app.go#L440-L448



#### Register Codec

`EncodingConfig`结构体是`app.go`文件的最后一个重要部分。该结构体的目标是定义将在整个应用程序中使用的编解码器。

https://github.com/cosmos/cosmos-sdk/blob/d9175200920e96bfa4182b5c8bc46d91b17a28a1/simapp/params/encoding.go#L9-L16


```go
// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Marshaler
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}
```


以下是四个字段各自含义的描述：

- `InterfaceRegistry`：Protobuf编解码器使用`InterfaceRegistry`处理接口，这些接口使用`google.protobuf.Any`进行了编码和解码（也称为“解压缩”）。以将任何内容视为包含`type_url`（接口的具体类型）和值（其编码字节）的结构.`InterfaceRegistry`提供了一种注册接口和实现的机制，可以从`Any`中安全地解压缩该接口和实现。应用程序的每个模块都实现`RegisterInterfaces`方法，该方法可用于注册模块自己的接口和实现。
- `Marshaler`: `Marshaler`是整个SDK使用的默认编解码器。它由用于对状态进行编码和解码的`BinaryMarshaler`和用于向用户输出数据（例如，在CLI中）的`JSONMarshaler`组成。默认情况下，SDK使用`Protobuf`作为`Marshaler`。
- `TxConfig`: `TxConfig`定义了客户端可以用来生成应用程序定义的具体事务类型的接口。目前，SDK处理两种事务类型：`SIGN_MODE_DIRECT`（使用Protobuf二进制代码进行有线编码）和`SIGN_MODE_LEGACY_AMINO_JSON`（依赖`Amino`）。[在此处详细了解交易](https://github.com/cosmos/cosmos-sdk/blob/master/docs/core/transactions.md)。
- `Amino`: SDK的某些旧版部分仍将`Amino`用于向后兼容。每个模块都暴露一个`RegisterLegacyAmino`方法，以在`Amino`中注册模块的特定类型。应用程序开发人员不应再使用此Amino编解码器，因为将在以后的版本中移除。


SDK公开了用于创建`EncodingConfig`的`MakeCodecs`函数。它使用Protobuf作为默认的`Marshaler`，并将其向下传递到应用程序的`appCodec`字段。它还会在应用程序的`legacyAmino`字段中实例化旧版`Amino`编解码器。

查看来自simapp的MakeCodecs的示例：

https://github.com/cosmos/cosmos-sdk/blob/d9175200920e96bfa4182b5c8bc46d91b17a28a1/simapp/app.go#L429-L435

```go
// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// simapp. It is useful for tests and clients who do not want to construct the
// full simapp
func MakeCodecs() (codec.Marshaler, *codec.LegacyAmino) {
	config := MakeEncodingConfig()
	return config.Marshaler, config.Amino
}
```


### Modules

`Modules`是SDK应用程序的灵魂。它们可以被视为状态机中的状态机。当交易通过ABCI从基础Tendermint引擎中继到应用程序时，它会被`baseapp`路由到适当的`modules`以便进行处理。这种范式使开发人员可以轻松构建复杂的状态机，因为他们所需的大多数模块通常已经存在。对于开发人员而言，构建SDK应用程序所涉及的大部分工作都是围绕构建其应用程序所需的尚不存在的自定义模块进行的，并将它们与已经存在于一个统一应用程序中的模块集成在一起。在应用程序目录中，标准做法是将模块存储在`x/`目录中（不要与SDK的`x/`目录混淆，该目录包含已构建的模块）。

#### Application Module Interface


Modules 必须实现Cosmos SDK的`AppModuleBasic`和`AppModule`中定义的接口。`AppModuleBasic`实现了模块的基本非依赖元素，例如编解码器，而`AppModule`则处理大部分模块方法（包括需要引用其他模块`keeper`的方法）。`AppModule`和`AppModuleBasic`类型都在名为`./module.go`的文件中定义。

`AppModule`在模块上公开了有用的方法的集合，这些方法有助于将模块组合成一个协调的应用程序。这些方法是从`module manager`中调用的，该模块管理应用程序的模块集合。


#### Message Types

`Message`是由实现消息接口的每个模块定义的对象。每个交易包含一个或多个`message`。

当全节点接收到有效的交易块时，Tendermint会通`DeliverTx`将每个交易都中继到应用程序。然后，应用程序处理交易：

- 1. 收到交易后，应用程序首先从`[]byte`反序列化出交易。
- 2. 然后，在提取交易中包含的消息之前，它会验证有关交易的一些事项，例如费用支付和签名。
- 3. 使用消息的`Type()`方法，`baseapp`可以将其路由到适当的模块的`handler`，以便对其进行处理。
- 4. 如果成功处理了消息，则状态将更新。

有关交易生命周期的更多详细信息，请单击[此处](https://github.com/cosmos/cosmos-sdk/blob/master/docs/basics/tx-lifecycle.md)。

模块开发人员在构建自己的模块时会创建自定义消息类型。通常的做法是在消息的类型声明之前加上`Msg`。例如，消息类型`MsgSend`允许用户传输代币：

https://github.com/cosmos/cosmos-sdk/blob/d9175200920e96bfa4182b5c8bc46d91b17a28a1/proto/cosmos/bank/v1beta1/tx.proto#L10-L19


```go
// MsgSend represents a message to send coins from one account to another.
message MsgSend {
  option (gogoproto.equal)           = false;
  option (gogoproto.goproto_getters) = false;

  string   from_address                    = 1 [(gogoproto.moretags) = "yaml:\"from_address\""];
  string   to_address                      = 2 [(gogoproto.moretags) = "yaml:\"to_address\""];
  repeated cosmos.base.v1beta1.Coin amount = 3
      [(gogoproto.nullable) = false, (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"];
}
```

它由`bank`模块的`handler`处理，该模块最终调用`auth`模块的`keeper`更新状态。


#### Handler


`handler`是指模块的一部分，负责处理由`baseapp`路由消息过来的消息。仅当通过`DeliverTx` ABCI消息从Tendermint中继事务时，才执行模块的`handler`。如果通过`CheckTx`中继交易，则仅执行*无状态*检查和与费用相关的有状态检查。为了更好地理解DeliverTx和CheckTx之间的区别以及有状态和无状态检查之间的区别，请单击[此处](https://github.com/cosmos/cosmos-sdk/blob/master/docs/basics/tx-lifecycle.md)。


模块的`handler`通常在一个名为`handler.go`的文件中定义，并包括：

- 开关函数`NewHandler`，用于将消息路由到适当的`handler`。并在`AppModule`中注册，以在应用程序的模块管理器中用于初始化应用程序的`router`。下面是来自`nameservice`教程的此类切换示例


- 模块定义的每种消息类型都有一个处理函数。开发人员在这些函数中编写消息处理逻辑。通常，这涉及进行状态检查以确保消息有效，并调用管理员的方法来更新状态。

处理程序函数返回类型为`sdk.Result`的结果，该结果通知应用程序消息是否已成功处理：
https://github.com/cosmos/cosmos-sdk/blob/d9175200920e96bfa4182b5c8bc46d91b17a28a1/types/result.go#L15-L40



#### gRPC Query Services


v0.40 Stargate版本中引入了gRPC查询服务。它们允许用户使用gRPC查询状态。它们是默认启用的，可以在`app.toml`中的`grpc.enable`和`grpc.address`字段下进行配置。

gRPC查询服务在模块的Protobuf定义中定义，特别是在query.proto内部。`query.proto`定义文件公开单个Query Protobuf服务。每个gRPC查询端点都与Query服务内部以rpc关键字开头的服务方法相对应。


Protobuf为每个模块生成一个`QueryServer`接口，其中包含所有服务方法。然后，模块的`keeper`需要通过提供每种服务方法的具体实现来实现此`QueryServer`接口。此具体实现是相应gRPC查询端点的处理程序。


最后，每个模块还应将`RegisterQueryService`方法实现为`AppModule`接口的一部分。此方法应调用生成的Protobuf代码提供的`RegisterQueryServer`函数。


#### Legacy Querier

传统查询器是在SDK中引入Protobuf和gRPC之前使用的查询器。它们适用于现有模块，但在以后的SDK版本中将不推荐使用。如果要开发新模块，则应首选gRPC查询服务，并且如果希望使用旧查询器，则只需实现`LegacyQuerierHandler`接口。

`Legacy queriers`与 `handler`非常相似，不同之处在于，传统查询器为用户查询状态而不是处理交易。终端用户从界面启动查询，该用户提供了`queryRoute`和一些数据。然后，通过使用`queryRoute`的`baseapp`的`handleQueryCustom`方法将`queryRoute`到正确的应用程序的查询器：

https://github.com/cosmos/cosmos-sdk/blob/d9175200920e96bfa4182b5c8bc46d91b17a28a1/baseapp/abci.go#L388-L418


模块的`Querier`在名为`keeper/querier.go`的文件中定义，包括：

- switch函数`NewQuerier`，用于将查询路由到适当的查询器功能。此函数返回查询器函数，并在`AppModule`中注册，以在应用程序的模块管理器中使用以初始化应用程序的查询`router`。

```go
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryResolve:
			return queryResolve(ctx, path[1:], req, keeper)
		case QueryWhois:
			return queryWhois(ctx, path[1:], req, keeper)
		case QueryNames:
			return queryNames(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}
```

- 需要由模块定义的每种数据类型的一个查询器功能。开发人员在这些函数中编写查询处理逻辑。通常，这涉及调用`keeper`的方法来查询状态并将其序列化为JSON。


#### Keeper


`Keepers`是其模块存储的看门人。要在模块的存储区中进行读取或写入，必须使用模块的`keeper`的方法。这由Cosmos SDK的对象功能模型( object-capabilities)来确保。只有持有store key的对象才能访问它，并且只有模块的`keeper`才应该持有该模块的store的key。

`Keepers` 通常在名为`keeper.go`的文件中定义。它包含`keeper`的类型定义和方法。


`keeper`类型定义通常包括：

- `multistore`中模块存储的键。
- 其他模块的`keeper`的引用。仅当`keeper`需要访问其他模块的存储（从它们读取或写入）时才需要。

- 应用程序的编解码器`codec`. `keeper`需要它在存储结构之前将其序列化，或在检索它们时将其反序列化，因为存储仅接受`[]byte`作为值。


与类型定义一起，`keeper.go`文件的下一个重要组成部分是`keeper`的构造函数`NewKeeper`。此函数使用`codec`，存储`keys`以及可能引用其他模块的`keeper`作为参数实例化上述类型的新`keeper`。从应用程序的构造函数中调用`NewKeeper`函数。文件的其余部分定义了`keeper`的方法，主要是`getter`和`setter`。



#### Command-Line, gRPC Services and REST Interfaces

每个模块都定义了命令行命令，gRPC服务和REST路由，以通过应用程序的界面向用户公开。这使最终用户可以创建模块中定义的类型的消息，或查询模块管理的状态的子集。

#### CLI

通常，与模块相关的命令在模块文件夹中名`client/cli`的文件夹中定义。CLI将命令分为两类，交易和查询，分别在`client/cli/tx.go`和`client/cli/query.go`中定义。两个命令行都是基于`Cobra`库构建:

- `Transactions`命令使用户可以生成新的交易，以便可以将它们包含在块中并最终更新状态。应该为模块中定义的每种消息类型创建一个命令。该命令使用用户提供的参数调用消息的构造函数，并将其包装到交易中。SDK处理签名和其他交易元数据的添加。
- `Queries`使用户可以查询模块定义的状态子集。查询命令将查询转发到应用程序的查询路由器，然后将查询路由到提供的`queryRoute`参数的适当查询器。


#### gRPC


gRPC是具有多种语言支持的现代开源高性能RPC框架。推荐使用gRPC与外部客户端（例如钱包，浏览器和其他后端服务）与节点进行交互的。

每个模块都可以公开gRPC端点（称为服务方法），并在模块的Protobuf `query.proto`文件中定义。服务方法由其名称，输入参数和输出响应定义。然后，该模块需要：

- 在`AppModuleBasic`上定义`RegisterGRPCRoutes`方法，以将客户端gRPC请求连接到模块内的正确handler。

- 对于每个服务方法，定义一个相应的handler。handler实现了服务gRPC请求所需的核心逻辑，并且位于`keeper/grpc_query.go`文件中。



#### gRPC-gateway REST Endpoints

某些外部客户端可能不希望使用gRPC。在这种情况下，SDK提供了gRPC网关服务，该服务将每个gRPC服务公开为相应的REST端点。请参阅grpc-gateway文档以了解更多信息。

REST endpoins 是使用Protobuf注释在Protobuf文件以及gRPC服务中定义的。想要公开REST查询的模块应在其rpc方法中添加`google.api.http`批注。默认情况下，SDK中定义的所有REST端点都有一个以`/cosmos/`前缀开头的URL。

SDK还提供了一个开发endpoint，可以为这些REST端点生成Swagger定义文件。可以在`api.swagger`键下的`app.toml`配置文件中启用此endpoints。




#### Legacy API REST Endpoints


该模块的传统REST接口使用户可以生成交易并通过对应用程序的传统API服务的REST调用查询状态。 REST路由在文件`client/rest/rest.go`中定义，该文件包含：


- `RegisterRoutes`函数，该函数注册文件中定义的每个路由。从应用程序内部使用的每个模块的主应用程序界面调用此函数。 SDK中使用的`router`是Gorilla's mux。
- 需要公开的每个查询或事务创建功能的自定义请求类型定义。这些自定义请求类型基于Cosmos SDK的基本请求类型：

https://github.com/cosmos/cosmos-sdk/blob/d9175200920e96bfa4182b5c8bc46d91b17a28a1/types/rest/rest.go#L62-L76

```go

// BaseReq defines a structure that can be embedded in other request structures
// that all share common "base" fields.
type BaseReq struct {
	From          string       `json:"from"`
	Memo          string       `json:"memo"`
	ChainID       string       `json:"chain_id"`
	AccountNumber uint64       `json:"account_number"`
	Sequence      uint64       `json:"sequence"`
	TimeoutHeight uint64       `json:"timeout_height"`
	Fees          sdk.Coins    `json:"fees"`
	GasPrices     sdk.DecCoins `json:"gas_prices"`
	Gas           string       `json:"gas"`
	GasAdjustment string       `json:"gas_adjustment"`
	Simulate      bool         `json:"simulate"`
}

```

- 每个请求的一个handler函数可以路由到给定的模块。这些功能实现了服务请求所需的核心逻辑。


这些旧版API endpoints 出现在SDK中是为了向后兼容，在下一发行版中将删除它们。



### Application Interface

接口使用户可以与全节点客户端进行交互。也就是说用户可以从全节点查询数据，或者创建并发送要由全节点中继并最终包含在块中的新交易。

主要界面是命令行界面。通过聚合在应用程序使用的每个模块中定义的CLI命令，可以构建SDK应用程序的CLI。应用程序的CLI与守护程序（例如`appd`）相同，并在名为`appd/main.go`的文件中定义。该文件包含：

- `main()`函数，该函数执行以构建`appd`接口客户端。此函数准备每个命令，然后在构建它们之前将它们添加到`rootCmd`中。在`appd`的根部，该功能添加了通用命令，例如`status`，`keys`和`config`，`query`，`tx`命令和`rest-server`。

- 通过调用`queryCmd`函数来添加查询命令。此函数返回一个`Cobra`命令，其中包含在每个应用程序模块中定义的查询命令（从`main()`函数作为sdk.`ModuleClients`数组传递），以及其他一些较低级别的查询命令，例如块查询或验证程序查询。通过使用CLI的命令`appd query [query]`可以调用查询命令。
- 通过调用`txCmd`函数添加交易命令。与`queryCmd`类似，该函数返回一个`Cobra`命令，其中包含在每个应用程序模块中定义的`tx`命令，以及较低级别的tx命令，例如事务签名或广播。通过使用CLI的命令`appd tx [tx]`可以调用Tx命令。

请参阅nameservice教程中的应用程序主命令行文件示例

https://github.com/cosmos/sdk-tutorials/blob/86a27321cf89cc637581762e953d0c07f8c78ece/nameservice/cmd/nscli/main.go



### Dependencies and Makefile

本部分是可选的，因为开发人员可以自由选择其依赖项管理器和项目构建方法。也就是说，当前最常用的版本控制框架是`go.mod`。它确保在整个应用程序中使用的每个库都以正确的版本导入。请参阅`nameservice`教程中的示例：
https://github.com/cosmos/sdk-tutorials/blob/c6754a1e313eb1ed973c5c91dcc606f2fd288811/go.mod#L1-L18



为了生成应用程序，通常使用`Makefile`。 `Makefile`主要确保在构建应用程序的两个入口点appd和appd之前运行go.mod。请参阅nameservice教程中的Makefile示例


https://github.com/cosmos/sdk-tutorials/blob/86a27321cf89cc637581762e953d0c07f8c78ece/nameservice/Makefile




## Transaction Lifecycle

> https://github.com/cosmos/cosmos-sdk/blob/master/docs/basics/tx-lifecycle.md


本文档描述了从创建到提交状态更改的事务生命周期。交易定义在其他文档中进行了描述。该交易将称为Tx。

### Transaction Creation

命令行界面是主要的应用程序界面之一。可以通过用户从命令行输入以下格式的命令来创建交易Tx，在`[command]`中提供交易的类型，在`[args]`中提供参数，并在`[flags]`中提供gas price等配置：

```
[appname] tx [command] [args] [flags]
```





