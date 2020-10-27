package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/youngqqcn/blog/x/blog/types"
)

func (k Keeper) CreatePost(ctx sdk.Context, post types.Post) {

	store := ctx.KVStore(k.storeKey)
	key := []byte(types.PostPrefix + post.ID)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(post)
	store.Set(key, value)
}
