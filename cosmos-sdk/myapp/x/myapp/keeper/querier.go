package keeper

import (
  // this line is used by starport scaffolding # 1
	"github.com/username/myapp/x/myapp/types"
		
	
		
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new querier for myapp clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
    // this line is used by starport scaffolding # 2
		case types.QueryListStudent:
			return listStudent(ctx, k)
		case types.QueryGetStudent:
			return getStudent(ctx, path[1:], k)
		case types.QueryListPost:
			return listPost(ctx, k)
		case types.QueryGetPost:
			return getPost(ctx, path[1:], k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown myapp query endpoint")
		}
	}
}
