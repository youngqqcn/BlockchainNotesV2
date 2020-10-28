package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/youngqqcn/scavenge/x/scavenge/types"
)

// CreateScavenge creates a scavenge
func (k Keeper) CreateScavenge(ctx sdk.Context, scavenge types.Scavenge) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.ScavengePrefix + scavenge.SolutionHash)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(scavenge)
	store.Set(key, value)
}

// GetScavenge returns the scavenge information
func (k Keeper) GetScavenge(ctx sdk.Context, key string) (types.Scavenge, error) {
	store := ctx.KVStore(k.storeKey)
	var scavenge types.Scavenge
	byteKey := []byte(types.ScavengePrefix + key)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &scavenge)
	if err != nil {
		return scavenge, err
	}
	return scavenge, nil
}

// SetScavenge sets a scavenge
func (k Keeper) SetScavenge(ctx sdk.Context, scavenge types.Scavenge) {
	scavengeKey := scavenge.SolutionHash
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(scavenge)
	key := []byte(types.ScavengePrefix + scavengeKey)
	store.Set(key, bz)
}

// DeleteScavenge deletes a scavenge
func (k Keeper) DeleteScavenge(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(types.ScavengePrefix + key))
}

//
// Functions used by querier
//

func listScavenge(ctx sdk.Context, k Keeper) ([]byte, error) {
	var scavengeList []types.Scavenge
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.ScavengePrefix))
	for ; iterator.Valid(); iterator.Next() {
		var scavenge types.Scavenge
		k.cdc.MustUnmarshalBinaryLengthPrefixed(store.Get(iterator.Key()), &scavenge)
		scavengeList = append(scavengeList, scavenge)
	}
	res := codec.MustMarshalJSONIndent(k.cdc, scavengeList)
	return res, nil
}

func getScavenge(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	key := path[0]
	scavenge, err := k.GetScavenge(ctx, key)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, scavenge)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// Get creator of the item
func (k Keeper) GetScavengeOwner(ctx sdk.Context, key string) sdk.AccAddress {
	scavenge, err := k.GetScavenge(ctx, key)
	if err != nil {
		return nil
	}
	return scavenge.Creator
}

// Check if the key exists in the store
func (k Keeper) ScavengeExists(ctx sdk.Context, key string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(types.ScavengePrefix + key))
}



// reveal是没有存入数据库的,   这里直接列出 已经解决的 scavenge  
func listReveal(ctx sdk.Context, k Keeper) ([]byte, error) {
	var revealList []types.Scavenge
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.ScavengePrefix))
	for ; iterator.Valid(); iterator.Next() {
		var scavenge types.Scavenge
		k.cdc.MustUnmarshalBinaryLengthPrefixed(store.Get(iterator.Key()), &scavenge)
		if scavenge.Scavenger == nil {
			continue
		}
		revealList = append(revealList, scavenge)
	}
	res := codec.MustMarshalJSONIndent(k.cdc, revealList)
	return res, nil
}

func getReveal(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	key := path[0]
	scavenge, err := k.GetScavenge(ctx, key)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, scavenge)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
