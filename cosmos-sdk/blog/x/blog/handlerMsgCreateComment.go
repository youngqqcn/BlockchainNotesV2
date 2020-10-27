package blog

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/youngqqcn/blog/x/blog/types"
	"github.com/youngqqcn/blog/x/blog/keeper"
)

func handleMsgCreateComment(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateComment) (*sdk.Result, error) {
	var comment = types.Comment{
		Creator: msg.Creator,
		ID:      msg.ID,
    	Body: msg.Body,
    	PostID: msg.PostID,
	}
	k.CreateComment(ctx, comment)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
