package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/youngqqcn/nameservice/x/nameservice/types"
)

// CreateWhois creates a whois
// func (k Keeper) CreateWhois(ctx sdk.Context, whois types.Whois) {
// 	store := ctx.KVStore(k.storeKey)

// 	key := []byte(types.WhoisPrefix + whois.ID)
// 	value := k.cdc.MustMarshalBinaryLengthPrefixed(whois)
// 	store.Set(key, value)
// }

// GetWhois returns the whois information
func (k Keeper) GetWhois(ctx sdk.Context, key string) (types.Whois, error) {
	store := ctx.KVStore(k.storeKey)
	var whois types.Whois
	byteKey := []byte(types.WhoisPrefix + key)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &whois)
	if err != nil {
		return whois, err
	}
	return whois, nil
}

// SetWhois sets a whois
func (k Keeper) SetWhois(ctx sdk.Context, name string, whois types.Whois) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(whois)
	key := []byte(types.WhoisPrefix + name)
	store.Set(key, bz)
}

// DeleteWhois deletes a whois
func (k Keeper) DeleteWhois(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(types.WhoisPrefix + key))
}

// Check if the key exists in the store
func (k Keeper) Exists(ctx sdk.Context, key string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(types.WhoisPrefix + key))
}

func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	whois, _ := k.GetWhois(ctx, name)
	return whois.Value
}

func (k Keeper) SetName(ctx sdk.Context, name string, value string) error {
	whois, err := k.GetWhois(ctx, name)
	if err != nil {
		return err
	}
	whois.Value = value
	k.SetWhois(ctx, name, whois)
	return nil
}

func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	whois, err := k.GetWhois(ctx, name)
	if err != nil {
		return false
	}
	return !whois.Owner.Empty()
}

// Get creator of the item
func (k Keeper) GetOwner(ctx sdk.Context, key string) sdk.AccAddress {
	whois, err := k.GetWhois(ctx, key)
	if err != nil {
		return nil
	}
	return whois.Owner
}

func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) error {
	whois, err := k.GetWhois(ctx, name)
	if err != nil {
		return err
	}
	k.SetWhois(ctx, name, whois)
	return nil
}

func (k Keeper) GetPrice(ctx sdk.Context, name string) (sdk.Coins, error) {
	whois, err := k.GetWhois(ctx, name)
	if err != nil {
		return sdk.Coins{}, err
	}
	return whois.Price, nil
}

func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) error {
	whois, err := k.GetWhois(ctx, name)
	if err != nil {
		return err
	}
	whois.Price = price
	k.SetWhois(ctx, name, whois)
	return nil
}

func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}

func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.WhoisPrefix))
}

//
// Functions used by querier
//

func listWhois(ctx sdk.Context, k Keeper) ([]byte, error) {
	// var whoisList []types.Whois
	var queryWhois []types.QueryWhois
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.WhoisPrefix))
	for ; iterator.Valid(); iterator.Next() {
		var whois types.Whois
		k.cdc.MustUnmarshalBinaryLengthPrefixed(store.Get(iterator.Key()), &whois)

		var qwhois = types.QueryWhois{
			Owner: whois.Owner,
			Price: whois.Price,
			Value: whois.Value,
			Name:  string(iterator.Key()),
		}
		queryWhois = append(queryWhois, qwhois)
		// whoisList = append(whoisList, whois)
	}
	res := codec.MustMarshalJSONIndent(k.cdc, queryWhois)
	return res, nil
}

func getWhois(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	key := path[0]
	whois, err := k.GetWhois(ctx, key)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, whois)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func resolveName(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {

	name := path[0]
	if len(name) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name is valid")
	}

	value := keeper.ResolveName(ctx, name)
	rsp, err := codec.MarshalJSONIndent(keeper.cdc, types.QueryResResolve{Value: value})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return rsp, nil
}
