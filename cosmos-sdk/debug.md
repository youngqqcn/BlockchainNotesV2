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


关于 Baseapp的剖析, 可以参考[baseapp](./README.md#Baseapp)
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


// AddRoute adds a route path to the router with a given handler. The route must
// be alphanumeric.
func (rtr *Router) AddRoute(path string, h sdk.Handler) sdk.Router {
	if !isAlphaNumeric(path) {
		panic("route expressions can only contain alphanumeric characters")
	}
	if rtr.routes[path] != nil {
		panic(fmt.Sprintf("route %s has already been initialized", path))
	}

	rtr.routes[path] = h  // 为route添加handler
	return rtr
}


//  nameservice NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		// this line is used by starport scaffolding # 1
		case types.MsgSetName:
			return handleMsgSetName(ctx, k, msg)
		case types.MsgBuyName:
			return handleMsgBuyName(ctx, k, msg)
		case types.MsgDeleteName:
			return handleMsgDeleteName(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}


// nameservice 的 NewQuerierHandler
// NewQuerierHandler returns the nameservice module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.keeper) // 调用了keeper的NewQuerier
}


// NewQuerier creates a new querier for nameservice clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
        // this line is used by starport scaffolding # 2
        // 查询数据库
		case types.QueryListWhois:
			return listWhois(ctx, k)
		case types.QueryGetWhois:
			return getWhois(ctx, path[1:], k)
		case types.QueryResolveName:
			return resolveName(ctx, path[1:], k )
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown nameservice query endpoint")
		}
	}
}



// RegisterRoutes registers all module routes and module querier routes
func (m *Manager) RegisterRoutes(router sdk.Router, queryRouter sdk.QueryRouter) {
	for _, module := range m.Modules {
		if module.Route() != "" {

            // module.Route() 一般是模块的名称, 例如: nameservice
            // module.NewHandler是一个根据msg类型进行switch的函数
            //       即不同的msg调用不同的handler处理 
            // 如上
			router.AddRoute(module.Route(), module.NewHandler())
		}
		if module.QuerierRoute() != "" {


            // module.QuerierRoute 一般是模块的名称, 例如: nameservice
            // module.NewQuerierHandler是 模块的 AppModule实现的 NewQuerierHandler,
            // 如上
			queryRouter.AddRoute(module.QuerierRoute(), module.NewQuerierHandler())
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


##  跟踪全节点处理一笔交易的过程

大致流程:

tx --> Tendermint-->baseapp:abci.go:CheckTx --> baseapp:runTx-->anteHandler进行模拟-->runMsgs(并没有真正完整运行,只是走了一个过场)

所以CheckTx主要执行了  validateBasicTxMsgs 和 anteHandler(模拟执行)


然后 ReCheck再一次对pending的交易进行检查


baseapp:abci.go:DeliverTx-->baseapp:runTx-->baseapp-->runMsgs:遍历msgs, 根据msg获取路由, 根据路由获取handler-->handler处理msg



baseapp:abci.go

```go
// CheckTx implements the ABCI interface and executes a tx in CheckTx mode. In
// CheckTx mode, messages are not executed. This means messages are only validated
// and only the AnteHandler is executed. State is persisted to the BaseApp's
// internal CheckTx state if the AnteHandler passes. Otherwise, the ResponseCheckTx
// will contain releveant error information. Regardless of tx execution outcome,
// the ResponseCheckTx will contain relevant gas execution context.
func (app *BaseApp) CheckTx(req abci.RequestCheckTx) abci.ResponseCheckTx {

    // 反序列花交易字节
	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		return sdkerrors.ResponseCheckTx(err, 0, 0, app.trace)
	}

	var mode runTxMode

	switch {
	case req.Type == abci.CheckTxType_New:
		mode = runTxModeCheck

	case req.Type == abci.CheckTxType_Recheck:
		mode = runTxModeReCheck

	default:
		panic(fmt.Sprintf("unknown RequestCheckTx type: %s", req.Type))
	}

	gInfo, result, err := app.runTx(mode, req.Tx, tx)
	if err != nil {
		return sdkerrors.ResponseCheckTx(err, gInfo.GasWanted, gInfo.GasUsed, app.trace)
	}

	return abci.ResponseCheckTx{
		GasWanted: int64(gInfo.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
		Log:       result.Log,
		Data:      result.Data,
		Events:    result.Events.ToABCIEvents(),
	}
}


// DeliverTx implements the ABCI interface and executes a tx in DeliverTx mode.
// State only gets persisted if all messages are valid and get executed successfully.
// Otherwise, the ResponseDeliverTx will contain releveant error information.
// Regardless of tx execution outcome, the ResponseDeliverTx will contain relevant
// gas execution context.
func (app *BaseApp) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		return sdkerrors.ResponseDeliverTx(err, 0, 0, app.trace)
	}

    // 执行交易
	gInfo, result, err := app.runTx(runTxModeDeliver, req.Tx, tx)
	if err != nil {
		return sdkerrors.ResponseDeliverTx(err, gInfo.GasWanted, gInfo.GasUsed, app.trace)
	}

	return abci.ResponseDeliverTx{
		GasWanted: int64(gInfo.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
		Log:       result.Log,
		Data:      result.Data,
		Events:    result.Events.ToABCIEvents(),
	}
}

// runTx processes a transaction within a given execution mode, encoded transaction
// bytes, and the decoded transaction itself. All state transitions occur through
// a cached Context depending on the mode provided. State only gets persisted
// if all messages get executed successfully and the execution mode is DeliverTx.
// Note, gas execution info is always returned. A reference to a Result is
// returned if the tx does not run out of gas and if all the messages are valid
// and execute successfully. An error is returned otherwise.
func (app *BaseApp) runTx(mode runTxMode, txBytes []byte, tx sdk.Tx) (gInfo sdk.GasInfo, result *sdk.Result, err error) {
	// NOTE: GasWanted should be returned by the AnteHandler. GasUsed is
	// determined by the GasMeter. We need access to the context to get the gas
	// meter so we initialize upfront.
	var gasWanted uint64

	ctx := app.getContextForTx(mode, txBytes)
	ms := ctx.MultiStore()

	// only run the tx if there is block gas remaining
	if mode == runTxModeDeliver && ctx.BlockGasMeter().IsOutOfGas() {
		gInfo = sdk.GasInfo{GasUsed: ctx.BlockGasMeter().GasConsumed()}
		return gInfo, nil, sdkerrors.Wrap(sdkerrors.ErrOutOfGas, "no block gas left to run tx")
	}

	var startingGas uint64
	if mode == runTxModeDeliver {
        // 如果是DeliverTx过来的, 则需要计算gas
		startingGas = ctx.BlockGasMeter().GasConsumed()
	}

	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			// TODO: Use ErrOutOfGas instead of ErrorOutOfGas which would allow us
			// to keep the stracktrace.
			case sdk.ErrorOutOfGas:
				err = sdkerrors.Wrap(
					sdkerrors.ErrOutOfGas, fmt.Sprintf(
						"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
						rType.Descriptor, gasWanted, ctx.GasMeter().GasConsumed(),
					),
				)

			default:
				err = sdkerrors.Wrap(
					sdkerrors.ErrPanic, fmt.Sprintf(
						"recovered: %v\nstack:\n%v", r, string(debug.Stack()),
					),
				)
			}

			result = nil
		}

		gInfo = sdk.GasInfo{GasWanted: gasWanted, GasUsed: ctx.GasMeter().GasConsumed()}
	}()

	// If BlockGasMeter() panics it will be caught by the above recover and will
	// return an error - in any case BlockGasMeter will consume gas past the limit.
	//
	// NOTE: This must exist in a separate defer function for the above recovery
	// to recover from this one.
	defer func() {
		if mode == runTxModeDeliver {
			ctx.BlockGasMeter().ConsumeGas(
				ctx.GasMeter().GasConsumedToLimit(), "block gas meter",
			)

			if ctx.BlockGasMeter().GasConsumed() < startingGas {
				panic(sdk.ErrorGasOverflow{Descriptor: "tx gas summation"})
			}
		}
	}()

    // 获取消息
    msgs := tx.GetMsgs()
    // 调用每个msg的ValidateBasic消息,对msg进行无状态检查
	if err := validateBasicTxMsgs(msgs); err != nil {
		return sdk.GasInfo{}, nil, err
	}

	if app.anteHandler != nil {
		var anteCtx sdk.Context
		var msCache sdk.CacheMultiStore

		// Cache wrap context before AnteHandler call in case it aborts.
		// This is required for both CheckTx and DeliverTx.
		// Ref: https://github.com/cosmos/cosmos-sdk/issues/2772
		//
		// NOTE: Alternatively, we could require that AnteHandler ensures that
		// writes do not happen if aborted/failed.  This may have some
        // performance benefits, but it'll be more difficult to get right.
        
        // 获取当前状态,并复制一个副本, 用于模拟执行
		anteCtx, msCache = app.cacheTxContext(ctx, txBytes)
		anteCtx = anteCtx.WithEventManager(sdk.NewEventManager())
		newCtx, err := app.anteHandler(anteCtx, tx, mode == runTxModeSimulate)

		if !newCtx.IsZero() {
			// At this point, newCtx.MultiStore() is cache-wrapped, or something else
			// replaced by the AnteHandler. We want the original multistore, not one
			// which was cache-wrapped for the AnteHandler.
			//
			// Also, in the case of the tx aborting, we need to track gas consumed via
			// the instantiated gas meter in the AnteHandler, so we update the context
			// prior to returning.
			ctx = newCtx.WithMultiStore(ms)
		}

		// GasMeter expected to be set in AnteHandler
		gasWanted = ctx.GasMeter().Limit()

		if err != nil {
			return gInfo, nil, err
		}

		msCache.Write()
	}

	// Create a new Context based off of the existing Context with a cache-wrapped
	// MultiStore in case message processing fails. At this point, the MultiStore
	// is doubly cached-wrapped.
	runMsgCtx, msCache := app.cacheTxContext(ctx, txBytes)

	// Attempt to execute all messages and only update state if all messages pass
	// and we're in DeliverTx. Note, runMsgs will never return a reference to a
	// Result if any single message fails or does not have a registered Handler.
	result, err = app.runMsgs(runMsgCtx, msgs, mode)
	if err == nil && mode == runTxModeDeliver {
		msCache.Write()
	}

	return gInfo, result, err
}


// runMsgs iterates through a list of messages and executes them with the provided
// Context and execution mode. Messages will only be executed during simulation
// and DeliverTx. An error is returned if any single message fails or if a
// Handler does not exist for a given message route. Otherwise, a reference to a
// Result is returned. The caller must not commit state if an error is returned.
func (app *BaseApp) runMsgs(ctx sdk.Context, msgs []sdk.Msg, mode runTxMode) (*sdk.Result, error) {
	msgLogs := make(sdk.ABCIMessageLogs, 0, len(msgs))
	data := make([]byte, 0, len(msgs))
	events := sdk.EmptyEvents()

	// NOTE: GasWanted is determined by the AnteHandler and GasUsed by the GasMeter.
	for i, msg := range msgs {

        // 如果是从 CheckTx 过来的, 则跳过
		// skip actual execution for (Re)CheckTx mode
		if mode == runTxModeCheck || mode == runTxModeReCheck {
			break
		}

        // 根据消息获取路由, 一般是模块名, 例如: nameservice
        msgRoute := msg.Route()
        //
		handler := app.router.Route(ctx, msgRoute)
		if handler == nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message route: %s; message index: %d", msgRoute, i)
		}

		msgResult, err := handler(ctx, msg)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "failed to execute message; message index: %d", i)
		}

		msgEvents := sdk.Events{
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type())),
		}
		msgEvents = msgEvents.AppendEvents(msgResult.Events)

		// append message events, data and logs
		//
		// Note: Each message result's data must be length-prefixed in order to
		// separate each result.
		events = events.AppendEvents(msgEvents)
		data = append(data, msgResult.Data...)
		msgLogs = append(msgLogs, sdk.NewABCIMessageLog(uint16(i), msgResult.Log, msgEvents))
	}

	return &sdk.Result{
		Data:   data,
		Log:    strings.TrimSpace(msgLogs.String()),
		Events: events,
	}, nil
}


```





##  跟踪客户端(命令行和REST)发起一笔交易的过程

AddCommand(
    queryCmd,  // 查询相关
    txCmd,  // 交易相关
    serveCmd, // rest-server
    )

txCmd
    SendTxCmd // 发送交易




```go

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [from_key_or_address] [to_address] [amount]",
		Short: "Create and sign a send tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
            msg := types.NewMsgSend(cliCtx.GetFromAddress(), to, coins)
            
            // 根据msg生成交易,签名,广播
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}




// GenerateOrBroadcastMsgs creates a StdTx given a series of messages. If
// the provided context has generate-only enabled, the tx will only be printed
// to STDOUT in a fully offline manner. Otherwise, the tx will be signed and
// broadcasted.
func GenerateOrBroadcastMsgs(cliCtx context.CLIContext, txBldr authtypes.TxBuilder, msgs []sdk.Msg) error {
	if cliCtx.GenerateOnly {
		return PrintUnsignedStdTx(txBldr, cliCtx, msgs)
	}

	return CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
}

// CompleteAndBroadcastTxCLI implements a utility function that facilitates
// sending a series of messages in a signed transaction given a TxBuilder and a
// QueryContext. It ensures that the account exists, has a proper number and
// sequence set. In addition, it builds and signs a transaction with the
// supplied messages. Finally, it broadcasts the signed transaction to a node.
func CompleteAndBroadcastTxCLI(txBldr authtypes.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) error {
	txBldr, err := PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	fromName := cliCtx.GetFromName()

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err = EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return err
		}

		gasEst := GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	if cliCtx.Simulate {
		return nil
	}

	if !cliCtx.SkipConfirm {
		stdSignMsg, err := txBldr.BuildSignMsg(msgs)
		if err != nil {
			return err
		}

		var json []byte
		if viper.GetBool(flags.FlagIndentResponse) {
			json, err = cliCtx.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
			if err != nil {
				panic(err)
			}
		} else {
			json = cliCtx.Codec.MustMarshalJSON(stdSignMsg)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return err
		}
	}

	// build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(fromName, keys.DefaultKeyPass, msgs)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	return cliCtx.PrintOutput(res)
}
```



## 跟踪客户端(REST)查询交易的过程


main.go:AddCommand:
    lcd.ServeCommand: 启动RESTserver

```go
// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
  // this line is used by starport scaffolding # 2
}


// ServeCommand will start the application REST service as a blocking process. It
// takes a codec to create a RestServer object and a function to register all
// necessary routes.
func ServeCommand(cdc *codec.Codec, registerRoutesFn func(*RestServer)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rs := NewRestServer(cdc)

			registerRoutesFn(rs)
			rs.registerSwaggerUI()

			// Start the rest server and return error if one exists
			err = rs.Start(
				viper.GetString(flags.FlagListenAddr),
				viper.GetInt(flags.FlagMaxOpenConnections),
				uint(viper.GetInt(flags.FlagRPCReadTimeout)),
				uint(viper.GetInt(flags.FlagRPCWriteTimeout)),
				viper.GetBool(flags.FlagUnsafeCORS),
			)

			return err
		},
	}

	return flags.RegisterRestServerFlags(cmd)
}


// RegisterRoutes registers nameservice-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// this line is used by starport scaffolding # 1
	r.HandleFunc("/nameservice/whois", buyNameHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/nameservice/whois", listWhoisHandler(cliCtx, "nameservice")).Methods("GET")
	r.HandleFunc("/nameservice/whois/{key}", getWhoisHandler(cliCtx, "nameservice")).Methods("GET")
	r.HandleFunc("/nameservice/whois/resolve-name/{name}", resolveNameHandler(cliCtx, "nameservice")).Methods("GET")
	r.HandleFunc("/nameservice/whois", setNameHandler(cliCtx)).Methods("PUT")
	r.HandleFunc("/nameservice/whois", deleteNameHandler(cliCtx)).Methods("DELETE")
}


```
