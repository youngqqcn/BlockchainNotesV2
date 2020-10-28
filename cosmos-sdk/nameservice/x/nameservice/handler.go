package nameservice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/youngqqcn/nameservice/x/nameservice/keeper"
	"github.com/youngqqcn/nameservice/x/nameservice/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
    // this line is used by starport scaffolding # 1
		case types.MsgCreateName:
			return handleMsgCreateName(ctx, k, msg)
		case types.MsgSetName:
			return handleMsgSetName(ctx, k, msg)
		case types.MsgDeleteName:
			return handleMsgDeleteName(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
