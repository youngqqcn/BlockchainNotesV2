package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/youngqqcn/nameservice/x/nameservice/keeper"
	"github.com/youngqqcn/nameservice/x/nameservice/types"
)

func handleMsgSetName(ctx sdk.Context, k keeper.Keeper, msg types.MsgSetName) (*sdk.Result, error) {

	// 检查owner是否匹配
	owner := k.GetOwner(ctx, msg.Name)
	if !msg.Owner.Equals(owner) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect Owner")
	}

	// 设置
	k.SetName(ctx, msg.Name, msg.Value)

	ctx.EventManager().Events().AppendEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeSetName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeName, msg.Name),
			sdk.NewAttribute(types.AttributeOwner, msg.Owner.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
