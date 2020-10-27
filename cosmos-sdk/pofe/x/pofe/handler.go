package pofe

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/youngqqcn/pofe/x/pofe/keeper"
	"github.com/youngqqcn/pofe/x/pofe/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
    // this line is used by starport scaffolding # 1
		case types.MsgCreateClaim:
			return handleMsgCreateClaim(ctx, k, msg)
		case types.MsgSetClaim:
			return handleMsgSetClaim(ctx, k, msg)
		case types.MsgDeleteClaim:
			return handleMsgDeleteClaim(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
