package blog

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/youngqqcn/blog/x/blog/types"
	"github.com/youngqqcn/blog/x/blog/keeper"
)

// Handle a message to delete name
func handleMsgDeleteComment(ctx sdk.Context, k keeper.Keeper, msg types.MsgDeleteComment) (*sdk.Result, error) {
	if !k.CommentExists(ctx, msg.ID) {
		// replace with ErrKeyNotFound for 0.39+
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msg.ID)
	}
	if !msg.Creator.Equals(k.GetCommentOwner(ctx, msg.ID)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	k.DeleteComment(ctx, msg.ID)
	return &sdk.Result{}, nil
}
