package myapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/username/myapp/x/myapp/types"
	"github.com/username/myapp/x/myapp/keeper"
)

func handleMsgSetPost(ctx sdk.Context, k keeper.Keeper, msg types.MsgSetPost) (*sdk.Result, error) {
	var post = types.Post{
		Creator: msg.Creator,
		ID:      msg.ID,
    	Title: msg.Title,
    	Body: msg.Body,
	}
	if !msg.Creator.Equals(k.GetPostOwner(ctx, msg.ID)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}

	k.SetPost(ctx, post)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
