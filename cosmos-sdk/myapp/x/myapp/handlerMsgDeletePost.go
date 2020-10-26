package myapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/username/myapp/x/myapp/types"
	"github.com/username/myapp/x/myapp/keeper"
)

// Handle a message to delete name
func handleMsgDeletePost(ctx sdk.Context, k keeper.Keeper, msg types.MsgDeletePost) (*sdk.Result, error) {
	if !k.PostExists(ctx, msg.ID) {
		// replace with ErrKeyNotFound for 0.39+
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msg.ID)
	}
	if !msg.Creator.Equals(k.GetPostOwner(ctx, msg.ID)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	k.DeletePost(ctx, msg.ID)
	return &sdk.Result{}, nil
}
