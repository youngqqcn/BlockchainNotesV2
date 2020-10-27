package blog

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/youngqqcn/blog/x/blog/types"
	"github.com/youngqqcn/blog/x/blog/keeper"
)

func handleMsgSetComment(ctx sdk.Context, k keeper.Keeper, msg types.MsgSetComment) (*sdk.Result, error) {
	var comment = types.Comment{
		Creator: msg.Creator,
		ID:      msg.ID,
    	Body: msg.Body,
    	PostID: msg.PostID,
	}
	if !msg.Creator.Equals(k.GetCommentOwner(ctx, msg.ID)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}

	k.SetComment(ctx, comment)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
