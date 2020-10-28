package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/youngqqcn/scavenge/x/scavenge/types"
)

// CreateCommit creates a commit
func (k Keeper) CreateCommit(ctx sdk.Context, commit types.Commit) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.CommitPrefix + commit.SolutionScavengerHash)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(commit)
	store.Set(key, value)
}

// GetCommit returns the commit information
func (k Keeper) GetCommit(ctx sdk.Context, key string) (types.Commit, error) {
	store := ctx.KVStore(k.storeKey)
	var commit types.Commit
	byteKey := []byte(types.CommitPrefix + key)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &commit)
	if err != nil {
		return commit, err
	}
	return commit, nil
}

//
// Functions used by querier
//

func listCommit(ctx sdk.Context, k Keeper) ([]byte, error) {
	var commitList []types.Commit
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CommitPrefix))
	for ; iterator.Valid(); iterator.Next() {
		var commit types.Commit
		k.cdc.MustUnmarshalBinaryLengthPrefixed(store.Get(iterator.Key()), &commit)
		commitList = append(commitList, commit)
	}
	res := codec.MustMarshalJSONIndent(k.cdc, commitList)
	return res, nil
}

func getCommit(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	key := path[0]
	commit, err := k.GetCommit(ctx, key)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, commit)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// // Get creator of the item
// func (k Keeper) GetCommitOwner(ctx sdk.Context, key string) sdk.AccAddress {
// 	commit, err := k.GetCommit(ctx, key)
// 	if err != nil {
// 		return nil
// 	}
// 	return commit.Creator
// }

// // Check if the key exists in the store
// func (k Keeper) CommitExists(ctx sdk.Context, key string) bool {
// 	store := ctx.KVStore(k.storeKey)
// 	return store.Has([]byte(types.CommitPrefix + key))
// }
