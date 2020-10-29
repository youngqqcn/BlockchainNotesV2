package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/youngqqcn/nameservice/x/nameservice/types"
)

// Keeper of the nameservice store
type Keeper struct {
	CoinKeeper bank.Keeper  // 用于访问 bank 模块
	storeKey   sdk.StoreKey // 存储相关
	cdc        *codec.Codec // 用于序列化和反序列化结构体
	// paramspace types.ParamSubspace
}

// NewKeeper creates a nameservice keeper
func NewKeeper(coinKeeper bank.Keeper, cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		CoinKeeper: coinKeeper,
		storeKey:   key,
		cdc:        cdc,
		// paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
