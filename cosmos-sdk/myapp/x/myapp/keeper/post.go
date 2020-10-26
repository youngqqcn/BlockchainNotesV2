package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/username/myapp/x/myapp/types"
    "github.com/cosmos/cosmos-sdk/codec"
)

// CreatePost creates a post
func (k Keeper) CreatePost(ctx sdk.Context, post types.Post) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.PostPrefix + post.ID)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(post)
	store.Set(key, value)
}

// GetPost returns the post information
func (k Keeper) GetPost(ctx sdk.Context, key string) (types.Post, error) {
	store := ctx.KVStore(k.storeKey)
	var post types.Post
	byteKey := []byte(types.PostPrefix + key)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &post)
	if err != nil {
		return post, err
	}
	return post, nil
}

// SetPost sets a post
func (k Keeper) SetPost(ctx sdk.Context, post types.Post) {
	postKey := post.ID
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(post)
	key := []byte(types.PostPrefix + postKey)
	store.Set(key, bz)
}

// DeletePost deletes a post
func (k Keeper) DeletePost(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(types.PostPrefix + key))
}

//
// Functions used by querier
//

func listPost(ctx sdk.Context, k Keeper) ([]byte, error) {
	var postList []types.Post
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.PostPrefix))
	for ; iterator.Valid(); iterator.Next() {
		var post types.Post
		k.cdc.MustUnmarshalBinaryLengthPrefixed(store.Get(iterator.Key()), &post)
		postList = append(postList, post)
	}
	res := codec.MustMarshalJSONIndent(k.cdc, postList)
	return res, nil
}

func getPost(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	key := path[0]
	post, err := k.GetPost(ctx, key)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, post)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// Get creator of the item
func (k Keeper) GetPostOwner(ctx sdk.Context, key string) sdk.AccAddress {
	post, err := k.GetPost(ctx, key)
	if err != nil {
		return nil
	}
	return post.Creator
}


// Check if the key exists in the store
func (k Keeper) PostExists(ctx sdk.Context, key string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(types.PostPrefix + key))
}
