package keeper

import (
	// this line is used by starport scaffolding # 1
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	blogtypes "github.com/youngqqcn/blog/x/blog/types"
)

// NewQuerier creates a new querier for blog clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		// this line is used by starport scaffolding # 2
		case blogtypes.QueryListPost:
			return listPost(ctx, k) // 获取所有post
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown blog query endpoint")
		}
	}
}
