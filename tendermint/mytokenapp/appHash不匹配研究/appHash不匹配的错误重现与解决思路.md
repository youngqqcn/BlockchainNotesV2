## 问题描述

问题:  某个验证节点更换了版本,导致对交易的处理结果不同(产生不同state),
    最终此验证节点与其他正常的节点的appHash不同
    
## 问题复现

复现操作:
- 1. node1, node2, node3, 正常出块 
- 2. 修改node4的源码, 并编译出mytokenapp_xx, 启动节点 (比如: 在transfer中对一个特定的value,做特殊处理, 
        以使node4处理结果与其他节点的state不同, 而出现appHash不同)
        
```go

func (app *MyTokenApp) transfer(fromAddress, toAddress crypto.Address, value int64) (bool, error) {

    // 节点4, 为了测试崩溃
    if value == 3571113 {
        value += 1
    }
    //....略

}
```
    
- 3.此时将node3停掉
- 4.先发行代币以保证测试账户有足够的余额, 向node4发送转账请求, 转账金额为设置的特定金额, 如例子中的:3571113, 
- 5.此时, node4会出现appHash 不匹配的错误(如下) , node1, node2, node3也会出现同样的错误:

``` 
version ==>: 89 hash ==> 1f56c78178ec4a8cb49ddf7ce18e67979f574508f38c87144860726c373c1d59
I[2020-11-06|14:04:45.360] Committed state                              module=state height=89 txs=1 appHash=1F56C78178EC4A8CB49DDF7CE18E67979F574508F38C87144860726C373C1D59
E[2020-11-06|14:04:46.758] enterPrevote: ProposalBlock is invalid       module=consensus height=90 round=0 err="wrong Block.Header.AppHash.  Expected 1F56C78178EC4A8CB49DDF7CE18E67979F574508F38C87144860726C373C1D59, got 7C55004B3DF63EDB4DFC4270C9F0164B6803A2EB76C2F72BC031721C2C2DD91B"
E[2020-11-06|14:04:49.532] enterPrevote: ProposalBlock is invalid       module=consensus height=90 round=1 err="wrong Block.Header.AppHash.  Expected 1F56C78178EC4A8CB49DDF7CE18E67979F574508F38C87144860726C373C1D59, got 7C55004B3DF63EDB4DFC4270C9F0164B6803A2EB76C2F72BC031721C2C2DD91B"

....

```

- 6.此时, 将node4停掉, 启动node3. node3的区块会跟上来,
- 7.经过一段时间的等待(大约五分钟, 为什么需要等待这么长时间?), node1, node2, node3 重新正常出块, 此时链恢复正常,  此时查询余额是 正常的余额, 而不是像node4多扣了1




## 问题分析


根据报错的信息, 定位到出错的代码是 `tendermint@v0.33.8/consensus/state.go`的`defaultDoPrevote`函数, 
此函数由 `enterPrevote` 调用, 总的调用链:

```

receiveRoutine -> cs.handleMsg -> tryAddVote -> addVote -> enterPrevote -> defaultDoPrevote



func (cs *State) defaultDoPrevote(height int64, round int) {
	logger := cs.Logger.With("height", height, "round", round)

	// If a block is locked, prevote that.
	if cs.LockedBlock != nil {
		logger.Info("enterPrevote: Block was locked")
		cs.signAddVote(types.PrevoteType, cs.LockedBlock.Hash(), cs.LockedBlockParts.Header())
		return
	}

	// If ProposalBlock is nil, prevote nil.
	if cs.ProposalBlock == nil {
		logger.Info("enterPrevote: ProposalBlock is nil")
		cs.signAddVote(types.PrevoteType, nil, types.PartSetHeader{})
		return
	}

	// Validate proposal block
    // 这里对区块进行验证, 进行各项验证, 包括对appHash进行验证
	err := cs.blockExec.ValidateBlock(cs.state, cs.ProposalBlock)
	if err != nil {
		// ProposalBlock is invalid, prevote nil.
		logger.Error("enterPrevote: ProposalBlock is invalid", "err", err)
		cs.signAddVote(types.PrevoteType, nil, types.PartSetHeader{})
		return
	}

	// Prevote cs.ProposalBlock
	// NOTE: the proposal signature is validated when it is received,
	// and the proposal block parts are validated as they are received (against the merkle hash in the proposal)
	logger.Info("enterPrevote: ProposalBlock is valid")
	cs.signAddVote(types.PrevoteType, cs.ProposalBlock.Hash(), cs.ProposalBlockParts.Header())
}
```


