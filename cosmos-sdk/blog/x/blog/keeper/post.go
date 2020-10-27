package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	blogtypes "github.com/youngqqcn/blog/x/blog/types"
)

// CreatePost 创建Post
func (k Keeper) CreatePost(ctx sdk.Context, post blogtypes.Post) {

	store := ctx.KVStore(k.storeKey)
	key := []byte(blogtypes.PostPrefix + post.ID)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(post)
	store.Set(key, value)
}

// 提供给 keeper/querier 使用
func listPost(ctx sdk.Context, k Keeper) ([]byte, error) {

	// 获取数据库中的所有的 post
	var postList []blogtypes.Post
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(blogtypes.PostPrefix))
	for ; iterator.Valid(); iterator.Next() {
		var post blogtypes.Post
		k.cdc.MustUnmarshalBinaryLengthPrefixed(store.Get(iterator.Key()), &post)
		postList = append(postList, post)
	}

	// 序列化
	res := codec.MustMarshalJSONIndent(k.cdc, postList)
	return res, nil
}
