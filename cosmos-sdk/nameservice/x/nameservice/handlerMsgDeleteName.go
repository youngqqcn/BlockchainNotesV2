package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/youngqqcn/nameservice/x/nameservice/types"
	"github.com/youngqqcn/nameservice/x/nameservice/keeper"
)

// Handle a message to delete name
func handleMsgDeleteName(ctx sdk.Context, k keeper.Keeper, msg types.MsgDeleteName) (*sdk.Result, error) {
	if !k.NameExists(ctx, msg.ID) {
		// replace with ErrKeyNotFound for 0.39+
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msg.ID)
	}
	if !msg.Creator.Equals(k.GetNameOwner(ctx, msg.ID)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	k.DeleteName(ctx, msg.ID)
	return &sdk.Result{}, nil
}
