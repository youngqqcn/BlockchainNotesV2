package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/youngqqcn/nameservice/x/nameservice/types"
    "github.com/cosmos/cosmos-sdk/codec"
)

// CreateName creates a name
func (k Keeper) CreateName(ctx sdk.Context, name types.Name) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.NamePrefix + name.ID)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(name)
	store.Set(key, value)
}

// GetName returns the name information
func (k Keeper) GetName(ctx sdk.Context, key string) (types.Name, error) {
	store := ctx.KVStore(k.storeKey)
	var name types.Name
	byteKey := []byte(types.NamePrefix + key)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &name)
	if err != nil {
		return name, err
	}
	return name, nil
}

// SetName sets a name
func (k Keeper) SetName(ctx sdk.Context, name types.Name) {
	nameKey := name.ID
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(name)
	key := []byte(types.NamePrefix + nameKey)
	store.Set(key, bz)
}

// DeleteName deletes a name
func (k Keeper) DeleteName(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(types.NamePrefix + key))
}

//
// Functions used by querier
//

func listName(ctx sdk.Context, k Keeper) ([]byte, error) {
	var nameList []types.Name
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.NamePrefix))
	for ; iterator.Valid(); iterator.Next() {
		var name types.Name
		k.cdc.MustUnmarshalBinaryLengthPrefixed(store.Get(iterator.Key()), &name)
		nameList = append(nameList, name)
	}
	res := codec.MustMarshalJSONIndent(k.cdc, nameList)
	return res, nil
}

func getName(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	key := path[0]
	name, err := k.GetName(ctx, key)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, name)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// Get creator of the item
func (k Keeper) GetNameOwner(ctx sdk.Context, key string) sdk.AccAddress {
	name, err := k.GetName(ctx, key)
	if err != nil {
		return nil
	}
	return name.Creator
}


// Check if the key exists in the store
func (k Keeper) NameExists(ctx sdk.Context, key string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(types.NamePrefix + key))
}
