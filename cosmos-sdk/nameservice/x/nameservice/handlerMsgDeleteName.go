package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/youngqqcn/nameservice/x/nameservice/keeper"
	"github.com/youngqqcn/nameservice/x/nameservice/types"
)

// Handle a message to delete name
func handleMsgDeleteName(ctx sdk.Context, k keeper.Keeper, msg types.MsgDeleteName) (*sdk.Result, error) {
	if !k.Exists(ctx, msg.Name) {
		// replace with ErrKeyNotFound for 0.39+
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msg.Name)
	}

	// 检查消息的发送方是否是域名的所有者
	if !msg.Owner.Equals(k.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	k.DeleteWhois(ctx, msg.Name)

	ctx.EventManager().Events().AppendEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeDeleteName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeName, msg.Name),
			sdk.NewAttribute(types.AttributeOwner, msg.Owner.String()),
		),
	)

	return &sdk.Result{}, nil
}
