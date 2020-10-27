package blog

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/youngqqcn/blog/x/blog/keeper"
	"github.com/youngqqcn/blog/x/blog/types"
)

func handlerMsgCreatePost(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreatePost) (*sdk.Result, error) {

	var post = types.Post{
		Creator: msg.Creator,
		ID:      msg.ID,
		Title:   msg.Title,
		Body:    msg.Body,
	}

	k.CreatePost(ctx, post) // keeper中的CreatePost,存储数据
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
