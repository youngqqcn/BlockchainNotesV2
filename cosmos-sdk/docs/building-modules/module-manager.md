<!--
order: 2
-->

# Module Manager

模块管理器

Cosmos SDK modules need to implement the [`AppModule` interfaces](#application-module-interfaces), in order to be managed by the application's [module manager](#module-manager). The module manager plays an important role in [`message` and `query` routing](../core/baseapp.md#routing), and allows application developers to set the order of execution of a variety of functions like [`BeginBlocker` and `EndBlocker`](../basics/app-anatomy.md#begingblocker-and-endblocker). {synopsis}



Cosmos SDK 模块需要实现`AppModule`接口, 以便能够被`module manager`管理. `module manager`在`message`和`query`的路由起很重要的作用, 同样允许应用开发者设置不同函数的启动顺序,像`BeginBlocker`和`EndBlocker`函数.


## Pre-requisite Readings

- [Introduction to SDK Modules](./intro.md) {prereq}

## Application Module Interfaces

Application module interfaces exist to facilitate(促进) the composition of modules together to form a functional SDK application. There are 3 main application module interfaces:

- [`AppModuleBasic`](#appmodulebasic) for independent module functionalities.
    > 独立的模块功能
- [`AppModule`](#appmodule) for inter-dependent module functionalities (except genesis-related functionalities).
    > 相互依赖的模块功能(除genesis相关的功能外)
- [`AppModuleGenesis`](#appmodulegenesis) for inter-dependent genesis-related module functionalities.
    > 相互依赖的genesis相关的模块功能

The `AppModuleBasic` interface exists to define independent methods of the module, i.e. those that do not depend on other modules in the application. This allows for the construction of the basic application structure early in the application definition, generally in the `init()` function of the [main application file](../basics/app-anatomy.md#core-application-file).

`AppModuleBasic`接口定义了模块的一些独立的方法.例如,它们不需要依赖其他应用程序模块.这样他们可以组成基础的早期应用程序结构, 通常在主应用程序的`init()`函数中.


The `AppModule` interface exists to define inter-dependent module methods. Many modules need to interract with other modules, typically through [`keeper`s](./keeper.md), which means there is a need for an interface where modules list their `keeper`s and other methods that require a reference to another module's object. `AppModule` interface also enables the module manager to set the order of execution between module's methods like `BeginBlock` and `EndBlock`, which is important in cases where the order of execution between modules matters in the context of the application.

`AppModule`接口定义一些相互依赖的模块方法. 很多模块需要与其他模块交互, 一般都是通过`keeper`s, 这也意味着模块需要提供一个方法用于列出他们的`keeper`,并且其他的方法需要引用其他模块对象. `AppModule`接口还允许`module manager`设置不同函数的启动顺序,像`BeginBlocker`和`EndBlocker`函数.


Lastly the interface for genesis functionality `AppModuleGenesis` is separated out from full module functionality `AppModule` so that modules which
are only used for genesis can take advantage of the `Module` patterns without having to define many placeholder functions.


最后,`AppModuleGenesis`接口用于genesis相关功能, 并且是独立于全节点功能的`AppModule`以便模块可以用于genesis, 也利用了`Module`模式的优点, 而不需要定义很多占位函数.



### `AppModuleBasic`

The `AppModuleBasic` interface defines the independent methods modules need to implement.

+++ https://github.com/cosmos/cosmos-sdk/blob/325be6ff215db457c6fc7668109640cd7fdac461/types/module/module.go#L49-L63




```go
// AppModuleBasic is the standard form for basic non-dependant elements of an application module.
type AppModuleBasic interface {

    // 返回模块的name
    Name() string
    
    // 为模块注册amino codec , 用于序列化和反序列化结构体, 以便于对数据结构进行持久化到KVStore
    RegisterLegacyAminoCodec(*codec.LegacyAmino)
    
    // 注册模块的接口类型 和 具体实现 proto.Message
	RegisterInterfaces(codectypes.InterfaceRegistry)

    // 返回一个默认的GenesisState 序列化 json.RawMessage. 默认GenesisState需要被开发者定义,主要用于测试
    DefaultGenesis(codec.JSONMarshaler) json.RawMessage
    
    // 对模块的GenesisState 进行验证返回 json.RawMessage 形式. 
	ValidateGenesis(codec.JSONMarshaler, client.TxEncodingConfig, json.RawMessage) error

    // client functionality
    // 注册REST 路由   
    RegisterRESTRoutes(client.Context, *mux.Router)
    
    // gRPC 路由
    RegisterGRPCRoutes(client.Context, *runtime.ServeMux)
    
    // 获取 root Tx 命令, 主要用于用户生成新的交易
    GetTxCmd() *cobra.Command
    
    // 返回 root query命令 , 用于命令行查询等操作
	GetQueryCmd() *cobra.Command
}
```

Let us go through the methods:

- `Name()`: Returns the name of the module as a `string`.
- `RegisterLegacyAminoCodec(*codec.LegacyAmino)`: Registers the `amino` codec for the module, which is used to marshal and unmarshal structs to/from `[]byte` in order to persist them in the module's `KVStore`.
- `RegisterInterfaces(codectypes.InterfaceRegistry)`: Registers a module's interface types and their concrete implementations as `proto.Message`.
- `DefaultGenesis(codec.JSONMarshaler)`: Returns a default [`GenesisState`](./genesis.md#genesisstate) for the module, marshalled to `json.RawMessage`. The default `GenesisState` need to be defined by the module developer and is primarily used for testing.
- `ValidateGenesis(codec.JSONMarshaler, client.TxEncodingConfig, json.RawMessage)`: Used to validate the `GenesisState` defined by a module, given in its `json.RawMessage` form. It will usually unmarshall the `json` before running a custom [`ValidateGenesis`](./genesis.md#validategenesis) function defined by the module developer.
- `RegisterRESTRoutes(client.Context, *mux.Router)`: Registers the REST routes for the module. These routes will be used to map REST request to the module in order to process them. See [rest.md](../interfaces/rest.md) for more.
- `RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux)`: Registers gRPC routes for the module. 
- `GetTxCmd()`: Returns the root [`Tx` command](./module-interfaces.md#tx) for the module. The subcommands of this root command are used by end-users to generate new transactions containing [`message`s](./messages-and-queries.md#queries) defined in the module.
- `GetQueryCmd()`: Return the root [`query` command](./module-interfaces.md#query) for the module. The subcommands of this root command are used by end-users to generate new queries to the subset of the state defined by the module.

All the `AppModuleBasic` of an application are managed by the [`BasicManager`](#basicmanager).

### `AppModuleGenesis`

The `AppModuleGenesis` interface is a simple embedding of the `AppModuleBasic` interface with two added methods.

+++ https://github.com/cosmos/cosmos-sdk/blob/325be6ff215db457c6fc7668109640cd7fdac461/types/module/module.go#L152-L158

```go
type AppModuleGenesis interface {
	AppModuleBasic

    // 初始化模块的状态子集, 在genesis时被调用
    InitGenesis(sdk.Context, codec.JSONMarshaler, json.RawMessage) []abci.ValidatorUpdate
    
    // 导出模块管理的最新的状态子集, 以用于genesis文件. 此方法用于:使用已经存在的链的状态启动一条新链
	ExportGenesis(sdk.Context, codec.JSONMarshaler) json.RawMessage
}
```


Let us go through the two added methods:

- `InitGenesis(sdk.Context, codec.JSONMarshaler, json.RawMessage)`: Initializes the subset of the state managed by the module. It is called at genesis (i.e. when the chain is first started).
- `ExportGenesis(sdk.Context, codec.JSONMarshaler)`: Exports the latest subset of the state managed by the module to be used in a new genesis file. `ExportGenesis` is called for each module when a new chain is started from the state of an existing chain.

It does not have its own manager, and exists separately from [`AppModule`](#appmodule) only for modules that exist only to implement genesis functionalities, so that they can be managed without having to implement all of `AppModule`'s methods. If the module is not only used during genesis, `InitGenesis(sdk.Context, codec.JSONMarshaler, json.RawMessage)` and `ExportGenesis(sdk.Context,  codec.JSONMarshaler)` will generally be defined as methods of the concrete type implementing the `AppModule` interface.

### `AppModule`

The `AppModule` interface defines the inter-dependent methods modules need to implement.

+++ https://github.com/cosmos/cosmos-sdk/blob/228728cce2af8d494c8b4e996d011492139b04ab/types/module/module.go#L160-L182

```go
// AppModule is the standard form for an application module
type AppModule interface {
	AppModuleGenesis

    // registers
    // 注册一些不变值, 如果不变值与预期值不同, 会触发相应的处理逻辑(通常是链挂掉)
	RegisterInvariants(sdk.InvariantRegistry)

	// routes
	Route() sdk.Route

	// Deprecated: use RegisterQueryService
	QuerierRoute() string

	// Deprecated: use RegisterQueryService
	LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier

    // RegisterServices allows a module to register services
    // 注册 gRPC 查询服务
	RegisterServices(Configurator)

    // ABCI
    
    // 允许模块的开发者自定义一些操作, 在每个区块开始前, 此函数会自动触发; 如果不需要, 则实现为空函数即可
    BeginBlock(sdk.Context, abci.RequestBeginBlock)
    
    // 允许模块的开发者自定义一些操作, 在每个区块结束时, 此函数会自动触发; 如果不需要, 则实现为空函数即可
	EndBlock(sdk.Context, abci.RequestEndBlock) []abci.ValidatorUpdate
}
```


`AppModule`s are managed by the [module manager](#manager). This interface embeds the `AppModuleGenesis` interface so that the manager can access all the independent and genesis inter-dependent methods of the module. This means that a concrete type implementing the `AppModule` interface must either implement all the methods of `AppModuleGenesis` (and by extension `AppModuleBasic`), or include a concrete type that does as parameter.


`AppModule`内含了`AppModuleGenesis`, 所以需要还需实现 `AppModuleGenesis` 接口

Let us go through the methods of `AppModule`:

- `RegisterInvariants(sdk.InvariantRegistry)`: Registers the [`invariants`](./invariants.md) of the module. If the invariants deviates(偏离) from its predicted value, the [`InvariantRegistry`](./invariants.md#registry) triggers appropriate logic (most often the chain will be halted).
- `Route()`: Returns the route for [`message`s](./messages-and-queries.md#messages) to be routed to the module by [`baseapp`](../core/baseapp.md#message-routing).
- `QuerierRoute()` (deprecated): Returns the name of the module's query route, for [`queries`](./messages-and-queries.md#queries) to be routes to the module by [`baseapp`](../core/baseapp.md#query-routing).
- `LegacyQuerierHandler(*codec.LegacyAmino)` (deprecated): Returns a [`querier`](./query-services.md#legacy-queriers) given the query `path`, in order to process the `query`.
- `RegisterQueryService(grpc.Server)`: Allows a module to register a gRPC query service.
- `BeginBlock(sdk.Context, abci.RequestBeginBlock)`: This method gives module developers the option to implement logic that is automatically triggered at the beginning of each block. Implement empty if no logic needs to be triggered at the beginning of each block for this module.
- `EndBlock(sdk.Context, abci.RequestEndBlock)`: This method gives module developers the option to implement logic that is automatically triggered at the beginning of each block. This is also where the module can inform the underlying consensus engine of validator set changes (e.g. the `staking` module). Implement empty if no logic needs to be triggered at the beginning of each block for this module.

### Implementing the Application Module Interfaces

Typically, the various application module interfaces are implemented in a file called `module.go`, located in the module's folder (e.g. `./x/module/module.go`).

Almost every module need to implement the `AppModuleBasic` and `AppModule` interfaces. If the module is only used for genesis, it will implement `AppModuleGenesis` instead of `AppModule`. The concrete type that implements the interface can add parameters that are required for the implementation of the various methods of the interface. For example, the `Route()` function often calls a `NewHandler(k keeper)` function defined in [`handler.go`](./handler.md) and therefore needs to pass the module's [`keeper`](./keeper.md) as parameter.

```go
// example
type AppModule struct {
	AppModuleBasic
	keeper       Keeper
}
```

In the example above, you can see that the `AppModule` concrete type references an `AppModuleBasic`, and not an `AppModuleGenesis`. That is because `AppModuleGenesis` only needs to be implemented in modules that focus on genesis-related functionalities. In most modules, the concrete `AppModule` type will have a reference to an `AppModuleBasic` and implement the two added methods of `AppModuleGenesis` directly in the `AppModule` type.



通常 `AppModule`直接实现`InitGenesis` 和 `ExportGenesis` 而不是实现 `AppModuleGenesis` 接口.



If no parameter is required (which is often the case for `AppModuleBasic`), just declare an empty concrete type like so:

```go
type AppModuleBasic struct{}
```

## Module Managers

Module managers are used to manage collections of `AppModuleBasic` and `AppModule`.

### `BasicManager`

The `BasicManager` is a structure that lists all the `AppModuleBasic` of an application:

+++ https://github.com/cosmos/cosmos-sdk/blob/325be6ff215db457c6fc7668109640cd7fdac461/types/module/module.go#L65-L66

```go
// BasicManager is a collection of AppModuleBasic
type BasicManager map[string]AppModuleBasic

```


It implements the following methods:

- `NewBasicManager(modules ...AppModuleBasic)`: Constructor function. It takes a list of the application's `AppModuleBasic` and builds a new `BasicManager`. This function is generally called in the `init()` function of [`app.go`](../basics/app-anatomy.md#core-application-file) to quickly initialize the independent elements of the application's modules (click [here](https://github.com/cosmos/gaia/blob/master/app/app.go#L59-L74) to see an example).
- `RegisterLegacyAminoCodec(cdc *codec.LegacyAmino)`: Registers the [`codec.LegacyAmino`s](../core/encoding.md#amino) of each of the application's `AppModuleBasic`. This function is usually called early on in the [application's construction](../basics/app-anatomy.md#constructor).
- `RegisterInterfaces(registry codectypes.InterfaceRegistry)`: Registers interface types and implementations of each of the application's `AppModuleBasic`.
- `DefaultGenesis(cdc codec.JSONMarshaler)`: Provides default genesis information for modules in the application by calling the [`DefaultGenesis(cdc codec.JSONMarshaler)`](./genesis.md#defaultgenesis) function of each module. It is used to construct a default genesis file for the application.
- `ValidateGenesis(cdc codec.JSONMarshaler, txEncCfg client.TxEncodingConfig, genesis map[string]json.RawMessage)`: Validates the genesis information modules by calling the [`ValidateGenesis(codec.JSONMarshaler, client.TxEncodingConfig, json.RawMessage)`](./genesis.md#validategenesis) function of each module.
- `RegisterRESTRoutes(ctx client.Context, rtr *mux.Router)`: Registers REST routes for modules by calling the [`RegisterRESTRoutes`](./module-interfaces.md#register-routes) function of each module. This function is usually called function from the `main.go` function of the [application's command-line interface](../interfaces/cli.md).
- `RegisterGRPCGatewayRoutes(clientCtx client.Context, rtr *runtime.ServeMux)`: Registers gRPC routes for modules.
- `AddTxCommands(rootTxCmd *cobra.Command)`: Adds modules' transaction commands to the application's [`rootTxCommand`](../interfaces/cli.md#transaction-commands). This function is usually called function from the `main.go` function of the [application's command-line interface](../interfaces/cli.md).
- `AddQueryCommands(rootQueryCmd *cobra.Command)`: Adds modules' query commands to the application's [`rootQueryCommand`](../interfaces/cli.md#query-commands). This function is usually called function from the `main.go` function of the [application's command-line interface](../interfaces/cli.md).

### `Manager`

The `Manager` is a structure that holds all the `AppModule` of an application, and defines the order of execution between several key components of these modules:

+++ https://github.com/cosmos/cosmos-sdk/blob/325be6ff215db457c6fc7668109640cd7fdac461/types/module/module.go#L223-L231


直接看代码


```go
// Manager defines a module manager that provides the high level utility for managing and executing
// operations for a group of modules
type Manager struct {
	Modules            map[string]AppModule
	OrderInitGenesis   []string
	OrderExportGenesis []string
	OrderBeginBlockers []string
	OrderEndBlockers   []string
}


// SetOrderInitGenesis sets the order of init genesis calls
func (m *Manager) SetOrderInitGenesis(moduleNames ...string) {
	m.OrderInitGenesis = moduleNames
}

// SetOrderExportGenesis sets the order of export genesis calls
func (m *Manager) SetOrderExportGenesis(moduleNames ...string) {
	m.OrderExportGenesis = moduleNames
}

// SetOrderBeginBlockers sets the order of set begin-blocker calls
func (m *Manager) SetOrderBeginBlockers(moduleNames ...string) {
	m.OrderBeginBlockers = moduleNames
}

// SetOrderEndBlockers sets the order of set end-blocker calls
func (m *Manager) SetOrderEndBlockers(moduleNames ...string) {
	m.OrderEndBlockers = moduleNames
}

// RegisterInvariants registers all module routes and module querier routes
func (m *Manager) RegisterInvariants(ir sdk.InvariantRegistry) {
	for _, module := range m.Modules {
		module.RegisterInvariants(ir)
	}
}

// RegisterRoutes registers all module routes and module querier routes
func (m *Manager) RegisterRoutes(router sdk.Router, queryRouter sdk.QueryRouter, legacyQuerierCdc *codec.LegacyAmino) {

    // 为每个模块注册路由
	for _, module := range m.Modules {
		if !module.Route().Empty() {
			router.AddRoute(module.Route())
		}
		if module.QuerierRoute() != "" {
			queryRouter.AddRoute(module.QuerierRoute(), module.LegacyQuerierHandler(legacyQuerierCdc))
		}
	}
}

// RegisterQueryServices registers all module query services
func (m *Manager) RegisterQueryServices(grpcRouter grpc.Server) {
	for _, module := range m.Modules {
        // 为每个模块注册grpc服务
		module.RegisterQueryService(grpcRouter)
	}
}

// InitGenesis performs init genesis functionality for modules
func (m *Manager) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, genesisData map[string]json.RawMessage) abci.ResponseInitChain {
	var validatorUpdates []abci.ValidatorUpdate
	for _, moduleName := range m.OrderInitGenesis {
		if genesisData[moduleName] == nil {
			continue
		}

        // 调用每个模块的InitGenesis
		moduleValUpdates := m.Modules[moduleName].InitGenesis(ctx, cdc, genesisData[moduleName])

		// use these validator updates if provided, the module manager assumes
		// only one module will update the validator set
		if len(moduleValUpdates) > 0 {
			if len(validatorUpdates) > 0 {
				panic("validator InitGenesis updates already set by a previous module")
			}
			validatorUpdates = moduleValUpdates
		}
	}

	return abci.ResponseInitChain{
		Validators: validatorUpdates,
	}
}

// ExportGenesis performs export genesis functionality for modules
func (m *Manager) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) map[string]json.RawMessage {
	genesisData := make(map[string]json.RawMessage)
	for _, moduleName := range m.OrderExportGenesis {

        // 调用每个模块的 ExportGenesis
		genesisData[moduleName] = m.Modules[moduleName].ExportGenesis(ctx, cdc)
	}

	return genesisData
}

// BeginBlock performs begin block functionality for all modules. It creates a
// child context with an event manager to aggregate events emitted from all
// modules.
func (m *Manager) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	for _, moduleName := range m.OrderBeginBlockers {
        // 调用每个模块的 ExportGenesis
		m.Modules[moduleName].BeginBlock(ctx, req)
	}

	return abci.ResponseBeginBlock{
		Events: ctx.EventManager().ABCIEvents(),
	}
}

// EndBlock performs end block functionality for all modules. It creates a
// child context with an event manager to aggregate events emitted from all
// modules.
func (m *Manager) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	validatorUpdates := []abci.ValidatorUpdate{}

	for _, moduleName := range m.OrderEndBlockers {

        // 调用每个模块的 ExportBlock
		moduleValUpdates := m.Modules[moduleName].EndBlock(ctx, req)

		// use these validator updates if provided, the module manager assumes
		// only one module will update the validator set
		if len(moduleValUpdates) > 0 {
			if len(validatorUpdates) > 0 {
				panic("validator EndBlock updates already set by a previous module")
			}

			validatorUpdates = moduleValUpdates
		}
	}

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Events:           ctx.EventManager().ABCIEvents(),
	}
}
```

The module manager is used throughout the application whenever an action on a collection of modules is required. It implements the following methods:

- `NewManager(modules ...AppModule)`: Constructor function. It takes a list of the application's `AppModule`s and builds a new `Manager`. It is generally called from the application's main [constructor function](../basics/app-anatomy.md#constructor-function).
- `SetOrderInitGenesis(moduleNames ...string)`: Sets the order in which the [`InitGenesis`](./genesis.md#initgenesis) function of each module will be called when the application is first started. This function is generally called from the application's main [constructor function](../basics/app-anatomy.md#constructor-function).
- `SetOrderExportGenesis(moduleNames ...string)`: Sets the order in which the [`ExportGenesis`](./genesis.md#exportgenesis) function of each module will be called in case of an export. This function is generally called from the application's main [constructor function](../basics/app-anatomy.md#constructor-function).
- `SetOrderBeginBlockers(moduleNames ...string)`: Sets the order in which the `BeginBlock()` function of each module will be called at the beginning of each block. This function is generally called from the application's main [constructor function](../basics/app-anatomy.md#constructor-function).
- `SetOrderEndBlockers(moduleNames ...string)`: Sets the order in which the `EndBlock()` function of each module will be called at the beginning of each block. This function is generally called from the application's main [constructor function](../basics/app-anatomy.md#constructor-function).
- `RegisterInvariants(ir sdk.InvariantRegistry)`: Registers the [invariants](./invariants.md) of each module.
- `RegisterRoutes(router sdk.Router, queryRouter sdk.QueryRouter, legacyQuerierCdc *codec.LegacyAmino)`: Registers module routes to the application's `router`, in order to route [`message`s](./messages-and-queries.md#messages) to the appropriate [`handler`](./handler.md), and module query routes to the application's `queryRouter`, in order to route [`queries`](./messages-and-queries.md#queries) to the appropriate [`querier`](./query-services.md#legacy-queriers).
- `RegisterQueryServices(grpcRouter grpc.Server)`: Registers all module gRPC query services.
- `InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, genesisData map[string]json.RawMessage)`: Calls the [`InitGenesis`](./genesis.md#initgenesis) function of each module when the application is first started, in the order defined in `OrderInitGenesis`. Returns an `abci.ResponseInitChain` to the underlying consensus engine, which can contain validator updates.
- `ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler)`: Calls the [`ExportGenesis`](./genesis.md#exportgenesis) function of each module, in the order defined in `OrderExportGenesis`. The export constructs a genesis file from a previously existing state, and is mainly used when a hard-fork upgrade of the chain is required.
- `BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock)`: At the beginning of each block, this function is called from [`baseapp`](../core/baseapp.md#beginblock) and, in turn, calls the [`BeginBlock`](./beginblock-endblock.md) function of each module, in the order defined in `OrderBeginBlockers`. It creates a child [context](../core/context.md) with an event manager to aggregate [events](../core/events.md) emitted from all modules. The function returns an `abci.ResponseBeginBlock` which contains the aforementioned events.
- `EndBlock(ctx sdk.Context, req abci.RequestEndBlock)`: At the end of each block, this function is called from [`baseapp`](../core/baseapp.md#endblock) and, in turn, calls the [`EndBlock`](./beginblock-endblock.md) function of each module, in the order defined in `OrderEndBlockers`. It creates a child [context](../core/context.md) with an event manager to aggregate [events](../core/events.md) emitted from all modules. The function returns an `abci.ResponseEndBlock` which contains the aforementioned events, as well as validator set updates (if any).

Here's an example of a concrete integration within an application:

+++ https://github.com/cosmos/cosmos-sdk/blob/2323f1ac0e9a69a0da6b43693061036134193464/simapp/app.go#L315-L362



```go

// NOTE: Any module instantiated in the module manager that is later modified
// must be passed by reference here.
app.mm = module.NewManager(
    genutil.NewAppModule(
        app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
        encodingConfig.TxConfig,
    ),
    auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
    vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
    bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
    capability.NewAppModule(appCodec, *app.CapabilityKeeper),
    crisis.NewAppModule(&app.CrisisKeeper),
    gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
    mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
    slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
    distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
    staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
    upgrade.NewAppModule(app.UpgradeKeeper),
    evidence.NewAppModule(app.EvidenceKeeper),
    ibc.NewAppModule(app.IBCKeeper),
    params.NewAppModule(app.ParamsKeeper),
    transferModule,
)

// During begin block slashing happens after distr.BeginBlocker so that
// there is nothing left over in the validator fee pool, so as to keep the
// CanWithdrawInvariant invariant.
// NOTE: staking module is required if HistoricalEntries param > 0
app.mm.SetOrderBeginBlockers(
    upgradetypes.ModuleName, minttypes.ModuleName, distrtypes.ModuleName, slashingtypes.ModuleName,
    evidencetypes.ModuleName, stakingtypes.ModuleName, ibchost.ModuleName,
)
app.mm.SetOrderEndBlockers(crisistypes.ModuleName, govtypes.ModuleName, stakingtypes.ModuleName)

// NOTE: The genutils module must occur after staking so that pools are
// properly initialized with tokens from genesis accounts.
// NOTE: Capability module must occur first so that it can initialize any capabilities
// so that other modules that want to create or claim capabilities afterwards in InitChain
// can do so safely.
app.mm.SetOrderInitGenesis(
    capabilitytypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName, distrtypes.ModuleName, stakingtypes.ModuleName,
    slashingtypes.ModuleName, govtypes.ModuleName, minttypes.ModuleName, crisistypes.ModuleName,
    ibchost.ModuleName, genutiltypes.ModuleName, evidencetypes.ModuleName, ibctransfertypes.ModuleName,
)

app.mm.RegisterInvariants(&app.CrisisKeeper)
app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
app.mm.RegisterServices(module.NewConfigurator(app.GRPCQueryRouter()))
```

## Next {hide}

Learn more about [`message`s and `queries`](./messages-and-queries.md) {hide}
