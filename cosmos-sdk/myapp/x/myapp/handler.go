package myapp

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/username/myapp/x/myapp/keeper"
	"github.com/username/myapp/x/myapp/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
    // this line is used by starport scaffolding # 1
		case types.MsgCreateStudent:
			return handleMsgCreateStudent(ctx, k, msg)
		case types.MsgSetStudent:
			return handleMsgSetStudent(ctx, k, msg)
		case types.MsgDeleteStudent:
			return handleMsgDeleteStudent(ctx, k, msg)
		case types.MsgCreatePost:
			return handleMsgCreatePost(ctx, k, msg)
		case types.MsgSetPost:
			return handleMsgSetPost(ctx, k, msg)
		case types.MsgDeletePost:
			return handleMsgDeletePost(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
