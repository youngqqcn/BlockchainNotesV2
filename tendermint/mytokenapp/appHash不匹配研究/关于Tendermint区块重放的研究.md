

/tendermint@v0.33.8/consensus/replay.go

```go

func (h *Handshaker) Handshake(proxyApp proxy.AppConns) error {

	// Handshake is done via ABCI Info on the query conn.
    // 这里通过ABCI调用客户端(app)的Info, 如果时本地的(直接用go的), 
	// 则调用localClient的InfoSync , InfoSync内部调用了 app.Info
	res, err := proxyApp.Query().InfoSync(proxy.RequestInfo) 
	if err != nil {
		return fmt.Errorf("error calling Info: %v", err)
	}

	blockHeight := res.LastBlockHeight // ABCI中  app.Info中返回的高度
	if blockHeight < 0 {
		return fmt.Errorf("got a negative last block height (%d) from the app", blockHeight)
	}
	appHash := res.LastBlockAppHash

	h.logger.Info("ABCI Handshake App Info",
		"height", blockHeight,
		"hash", fmt.Sprintf("%X", appHash),
		"software-version", res.Version,
		"protocol-version", res.AppVersion,
	)

	// Set AppVersion on the state.
	if h.initialState.Version.Consensus.App != version.Protocol(res.AppVersion) {
		h.initialState.Version.Consensus.App = version.Protocol(res.AppVersion)
		sm.SaveState(h.stateDB, h.initialState)
	}

	// Replay blocks up to the latest in the blockstore.
    // 开始重放区块
	_, err = h.ReplayBlocks(h.initialState, appHash, blockHeight, proxyApp)
	if err != nil {
		return fmt.Errorf("error on replay: %v", err)
	}

	h.logger.Info("Completed ABCI Handshake - Tendermint and App are synced",
		"appHeight", blockHeight, "appHash", fmt.Sprintf("%X", appHash))

	// TODO: (on restart) replay mempool

	return nil
}

```

/home/yqq/go/pkg/mod/github.com/tendermint/tendermint@v0.33.8/abci/client/local_client.go

```go
func (app *localClient) InfoSync(req types.RequestInfo) (*types.ResponseInfo, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.Info(req)
	return &res, nil
}

```



```go

// ReplayBlocks replays all blocks since appBlockHeight and ensures the result
// matches the current state.
// Returns the final AppHash or an error.
func (h *Handshaker) ReplayBlocks(
	state sm.State,
	appHash []byte,
	appBlockHeight int64, // 这个是ABCI 即 app.Info()中返回的高度
	proxyApp proxy.AppConns,
) ([]byte, error) {
	storeBlockBase := h.store.Base()
	storeBlockHeight := h.store.Height() //这个是
	stateBlockHeight := state.LastBlockHeight
	h.logger.Info(
		"ABCI Replay Blocks",
		"appHeight",
		appBlockHeight,
		"storeHeight",
		storeBlockHeight,
		"stateHeight",
		stateBlockHeight)

	// If appBlockHeight == 0 it means that we are at genesis and hence should send InitChain.
    // 如果app.Info()返回的是 0, 则从高度0开始进行区块回放
	if appBlockHeight == 0 {
		validators := make([]*types.Validator, len(h.genDoc.Validators))
		for i, val := range h.genDoc.Validators {
			validators[i] = types.NewValidator(val.PubKey, val.Power)
		}
        
        // 获取初始化参数
		validatorSet := types.NewValidatorSet(validators)
		nextVals := types.TM2PB.ValidatorUpdates(validatorSet)
		csParams := types.TM2PB.ConsensusParams(h.genDoc.ConsensusParams)
		
        //请求初始化
        req := abci.RequestInitChain{
			Time:            h.genDoc.GenesisTime,
			ChainId:         h.genDoc.ChainID,
			ConsensusParams: csParams,
			Validators:      nextVals,
			AppStateBytes:   h.genDoc.AppState,
		}
        
        // 初始化
		res, err := proxyApp.Consensus().InitChainSync(req)
		if err != nil {
			return nil, err
		}
        
        // 如果当前链的状态真的是创世区块, 则只需要初始化一些参数即可
		if stateBlockHeight == 0 { //we only update state when we are in initial state
			// If the app returned validators or consensus params, update the state.
			if len(res.Validators) > 0 {
				vals, err := types.PB2TM.ValidatorUpdates(res.Validators)
				if err != nil {
					return nil, err
				}
				state.Validators = types.NewValidatorSet(vals)
				state.NextValidators = types.NewValidatorSet(vals)
			} else if len(h.genDoc.Validators) == 0 {
				// If validator set is not set in genesis and still empty after InitChain, exit.
				return nil, fmt.Errorf("validator set is nil in genesis and still empty after InitChain")
			}

			if res.ConsensusParams != nil {
				state.ConsensusParams = state.ConsensusParams.Update(res.ConsensusParams)
			}
			sm.SaveState(h.stateDB, state)
		}
	}

	// First handle edge cases and constraints on the storeBlockHeight and storeBlockBase.
    // 如果 app.Info() 返回的不是0, 则需要
	switch {
	case storeBlockHeight == 0:
		assertAppHashEqualsOneFromState(appHash, state)
		return appHash, nil

	case appBlockHeight < storeBlockBase-1:
		// the app is too far behind truncated store (can be 1 behind since we replay the next)
		return appHash, sm.ErrAppBlockHeightTooLow{AppHeight: appBlockHeight, StoreBase: storeBlockBase}

	case storeBlockHeight < appBlockHeight:
		// the app should never be ahead of the store (but this is under app's control)
		return appHash, sm.ErrAppBlockHeightTooHigh{CoreHeight: storeBlockHeight, AppHeight: appBlockHeight}

	case storeBlockHeight < stateBlockHeight:
		// the state should never be ahead of the store (this is under tendermint's control)
		panic(fmt.Sprintf("StateBlockHeight (%d) > StoreBlockHeight (%d)", stateBlockHeight, storeBlockHeight))

	case storeBlockHeight > stateBlockHeight+1:
		// store should be at most one ahead of the state (this is under tendermint's control)
		panic(fmt.Sprintf("StoreBlockHeight (%d) > StateBlockHeight + 1 (%d)", storeBlockHeight, stateBlockHeight+1))
	}

	var err error
	// Now either store is equal to state, or one ahead.
	// For each, consider all cases of where the app could be, given app <= store
	if storeBlockHeight == stateBlockHeight {
		// Tendermint ran Commit and saved the state.
		// Either the app is asking for replay, or we're all synced up.
		if appBlockHeight < storeBlockHeight {
			// the app is behind, so replay blocks, but no need to go through WAL (state is already synced to store)
			return h.replayBlocks(state, proxyApp, appBlockHeight, storeBlockHeight, false)

		} else if appBlockHeight == storeBlockHeight {
			// We're good!
			assertAppHashEqualsOneFromState(appHash, state)
			return appHash, nil
		}

	} else if storeBlockHeight == stateBlockHeight+1 {
		// We saved the block in the store but haven't updated the state,
		// so we'll need to replay a block using the WAL.
		switch {
		case appBlockHeight < stateBlockHeight:
			// the app is further behind than it should be, so replay blocks
			// but leave the last block to go through the WAL
			return h.replayBlocks(state, proxyApp, appBlockHeight, storeBlockHeight, true)

		case appBlockHeight == stateBlockHeight:
			// We haven't run Commit (both the state and app are one block behind),
			// so replayBlock with the real app.
			// NOTE: We could instead use the cs.WAL on cs.Start,
			// but we'd have to allow the WAL to replay a block that wrote it's #ENDHEIGHT
			h.logger.Info("Replay last block using real app")
			state, err = h.replayBlock(state, storeBlockHeight, proxyApp.Consensus())
			return state.AppHash, err

		case appBlockHeight == storeBlockHeight:
			// We ran Commit, but didn't save the state, so replayBlock with mock app.
			abciResponses, err := sm.LoadABCIResponses(h.stateDB, storeBlockHeight)
			if err != nil {
				return nil, err
			}
			mockApp := newMockProxyApp(appHash, abciResponses)
			h.logger.Info("Replay last block using mock app")
			state, err = h.replayBlock(state, storeBlockHeight, mockApp)
			return state.AppHash, err
		}

	}

	panic(fmt.Sprintf("uncovered case! appHeight: %d, storeHeight: %d, stateHeight: %d",
		appBlockHeight, storeBlockHeight, stateBlockHeight))
}

func (h *Handshaker) replayBlocks(
	state sm.State,
	proxyApp proxy.AppConns,
	appBlockHeight,
	storeBlockHeight int64,
	mutateState bool) ([]byte, error) {
	// App is further behind than it should be, so we need to replay blocks.
	// We replay all blocks from appBlockHeight+1.
	//
	// Note that we don't have an old version of the state,
	// so we by-pass state validation/mutation using sm.ExecCommitBlock.
	// This also means we won't be saving validator sets if they change during this period.
	// TODO: Load the historical information to fix this and just use state.ApplyBlock
	//
	// If mutateState == true, the final block is replayed with h.replayBlock()

	var appHash []byte
	var err error
	finalBlock := storeBlockHeight
	if mutateState {
		finalBlock--
	}
	for i := appBlockHeight + 1; i <= finalBlock; i++ {
		h.logger.Info("Applying block", "height", i)
		block := h.store.LoadBlock(i)
		// Extra check to ensure the app was not changed in a way it shouldn't have.
		if len(appHash) > 0 {
			assertAppHashEqualsOneFromBlock(appHash, block)
		}

		appHash, err = sm.ExecCommitBlock(proxyApp.Consensus(), block, h.logger, h.stateDB)
		if err != nil {
			return nil, err
		}

		h.nBlocks++
	}

	if mutateState {
		// sync the final block
		state, err = h.replayBlock(state, storeBlockHeight, proxyApp.Consensus())
		if err != nil {
			return nil, err
		}
		appHash = state.AppHash
	}

	assertAppHashEqualsOneFromState(appHash, state)
	return appHash, nil
}

// ApplyBlock on the proxyApp with the last block.
func (h *Handshaker) replayBlock(state sm.State, height int64, proxyApp proxy.AppConnConsensus) (sm.State, error) {
	block := h.store.LoadBlock(height)
	meta := h.store.LoadBlockMeta(height)

	blockExec := sm.NewBlockExecutor(h.stateDB, h.logger, proxyApp, mock.Mempool{}, sm.MockEvidencePool{})
	blockExec.SetEventBus(h.eventBus)

	var err error
	state, _, err = blockExec.ApplyBlock(state, meta.BlockID, block)
	if err != nil {
		return sm.State{}, err
	}

	h.nBlocks++

	return state, nil
}

```







        
    




























