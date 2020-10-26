package myapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/username/myapp/x/myapp/types"
	"github.com/username/myapp/x/myapp/keeper"
)

func handleMsgCreatePost(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreatePost) (*sdk.Result, error) {
	var post = types.Post{
		Creator: msg.Creator,
		ID:      msg.ID,
    	Title: msg.Title,
    	Body: msg.Body,
	}
	k.CreatePost(ctx, post)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
