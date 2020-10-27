package keeper

import (
	// this line is used by starport scaffolding # 1
	"github.com/youngqqcn/blog/x/blog/types"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new querier for blog clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		// this line is used by starport scaffolding # 2
		case types.QueryListComment:
			return listComment(ctx, k)
		case types.QueryGetComment:
			return getComment(ctx, path[1:], k)
		case types.QueryListPost:
			return listPost(ctx, k) // 获取所有post
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown blog query endpoint")
		}
	}
}
