# 调试跟踪Cosmos应用程序记录

- 调试目的: 深入理解cosmos-sdk官方文档, 理解 module, module manager, message, handler, keeper, store, 以及msg的处理流程
- 调试方式: 以`nameservice`示例对象,
  - 跟踪app的初始化过程
  - 跟踪全节点处理一笔交易的过程,
  - 跟踪客户端(命令行和REST)发起一笔交易的过程
  - 跟踪客户端(命令和和REST)查询交易的过程



## 跟踪app的初始化过程

app.go 创建 ModuleBasics 
app.go:MakeCodec() 注册codec

app/prefix.go:设置config 

main.go:newApp: 
1.创建CommitKVStoreCacheManager

```go
    CommitKVStoreCacheManager struct {
		cacheSize uint
		caches    map[string]types.CommitKVStore
	}
```


```go

return app.NewInitApp(
logger, db, traceStore, true, invCheckPeriod,
baseapp.SetPruning(pruningOpts),
// 设置最小gasPrice
baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
// 设置停止高度??
baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
// 设置停止时间?? 
baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
baseapp.SetInterBlockCache(cache),


```


app.go
```go

func NewBaseApp(
	name string, logger log.Logger, db dbm.DB, txDecoder sdk.TxDecoder, options ...func(*BaseApp),
) *BaseApp {

	app := &BaseApp{
		logger:         logger,
		name:           name,
		db:             db,
		cms:            store.NewCommitMultiStore(db),
		storeLoader:    DefaultStoreLoader,
		router:         NewRouter(),
		queryRouter:    NewQueryRouter(),
		txDecoder:      txDecoder,
		fauxMerkleMode: false,
		trace:          false,
    }
    
    // 执行所有选项函数
	for _, option := range options {
		option(app) 
	}

	if app.interBlockCache != nil {
		app.cms.SetInterBlockCache(app.interBlockCache)
	}

	return app
}


// NewKVStoreKeys returns a map of new  pointers to KVStoreKey's.
// Uses pointers so keys don't collide.
// 创建 KVStoreKey 用于管理KVStore
func NewKVStoreKeys(names ...string) map[string]*KVStoreKey {
	keys := make(map[string]*KVStoreKey)
	for _, name := range names {
		keys[name] = NewKVStoreKey(name)
	}
	return keys
}


// Keeper of the global paramstore
type Keeper struct {
	cdc    *codec.Codec
	key    sdk.StoreKey
	tkey   sdk.StoreKey
	spaces map[string]*Subspace
}

// NewKeeper constructs a params keeper
func NewKeeper(cdc *codec.Codec, key, tkey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:    cdc,
		key:    key,
		tkey:   tkey,
		spaces: make(map[string]*Subspace),
	}
}

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountKVStores(keys map[string]*sdk.KVStoreKey) {
	for _, key := range keys {
		if !app.fauxMerkleMode { // 需要进行merkle proof
			app.MountStore(key, sdk.StoreTypeIAVL)
		} else { // 不需要merkle proof
			// StoreTypeDB doesn't do anything upon commit, and it doesn't
            // retain history, but it's useful for faster simulation.
            // 不会保留历史记录
			app.MountStore(key, sdk.StoreTypeDB)
		}
	}
}

func NewInitApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *NewQnsApp {
	cdc := MakeCodec()

    // 创建 baseapp
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

    // 创建 KVStoreKey 用于管理KVStore
	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey,
		auth.StoreKey,
		staking.StoreKey,
		supply.StoreKey,
		params.StoreKey,
		nameservicetypes.StoreKey,
		// this line is used by starport scaffolding # 5
	)

    // 创建临时的StoreKeys?
    // TransientStoreKey is used for indexing transient stores in a MultiStore
    // TransientStoreKey 用于在multistore中索引临时的store,
	tKeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	var app = &NewQnsApp{
		BaseApp:        bApp, // baseapp
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		subspaces:      make(map[string]params.Subspace),
    }
    

    // 构造keeper

	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	// this line is used by starport scaffolding # 5.1

	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		app.subspaces[auth.ModuleName],
		auth.ProtoBaseAccount,
	)

	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.subspaces[bank.ModuleName],
		app.ModuleAccountAddrs(),
	)

	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		maccPerms,
	)

	stakingKeeper := staking.NewKeeper(
		app.cdc,
		keys[staking.StoreKey],
		app.supplyKeeper,
		app.subspaces[staking.ModuleName],
	)

	// this line is used by starport scaffolding # 5.2

	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
		// this line is used by starport scaffolding # 5.3
		),
	)

	app.nameserviceKeeper = nameservicekeeper.NewKeeper(
		app.bankKeeper,
		app.cdc,
		keys[nameservicetypes.StoreKey],
	)

	// this line is used by starport scaffolding # 4


    // 构造模块管理器
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		nameservice.NewAppModule(app.nameserviceKeeper, app.bankKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		// this line is used by starport scaffolding # 6
	)

    // 设置EndBlocker的执行顺序
	app.mm.SetOrderEndBlockers(
		staking.ModuleName,
		// this line is used by starport scaffolding # 6.1
	)

    // 设置InitGenesis的执行顺序
	app.mm.SetOrderInitGenesis(
		// this line is used by starport scaffolding # 6.2
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		nameservicetypes.ModuleName,
		supply.ModuleName,
		genutil.ModuleName,
		// this line is used by starport scaffolding # 7
	)

    // 注册模块路由, 和 查询路由(querierRoute)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

    // 设置ABCI的处理函数
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

    // 设置 AnteHandler 用于对msg进行检查(有状态,无状态), 例如,验证签名,检查gas等
	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

    //挂载所有数据库( 根据前面的 KeyStoreKeys)
    app.MountKVStores(keys)
    
    //挂载所有临时的数据库(根据临时的KeyStoreKeys)
	app.MountTransientStores(tKeys)

    //加载最新的主库状态(main)
	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}
```



