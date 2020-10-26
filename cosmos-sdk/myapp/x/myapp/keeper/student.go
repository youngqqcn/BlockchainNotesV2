package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/username/myapp/x/myapp/types"
    "github.com/cosmos/cosmos-sdk/codec"
)

// CreateStudent creates a student
func (k Keeper) CreateStudent(ctx sdk.Context, student types.Student) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.StudentPrefix + student.ID)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(student)
	store.Set(key, value)
}

// GetStudent returns the student information
func (k Keeper) GetStudent(ctx sdk.Context, key string) (types.Student, error) {
	store := ctx.KVStore(k.storeKey)
	var student types.Student
	byteKey := []byte(types.StudentPrefix + key)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &student)
	if err != nil {
		return student, err
	}
	return student, nil
}

// SetStudent sets a student
func (k Keeper) SetStudent(ctx sdk.Context, student types.Student) {
	studentKey := student.ID
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(student)
	key := []byte(types.StudentPrefix + studentKey)
	store.Set(key, bz)
}

// DeleteStudent deletes a student
func (k Keeper) DeleteStudent(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(types.StudentPrefix + key))
}

//
// Functions used by querier
//

func listStudent(ctx sdk.Context, k Keeper) ([]byte, error) {
	var studentList []types.Student
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.StudentPrefix))
	for ; iterator.Valid(); iterator.Next() {
		var student types.Student
		k.cdc.MustUnmarshalBinaryLengthPrefixed(store.Get(iterator.Key()), &student)
		studentList = append(studentList, student)
	}
	res := codec.MustMarshalJSONIndent(k.cdc, studentList)
	return res, nil
}

func getStudent(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	key := path[0]
	student, err := k.GetStudent(ctx, key)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, student)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// Get creator of the item
func (k Keeper) GetStudentOwner(ctx sdk.Context, key string) sdk.AccAddress {
	student, err := k.GetStudent(ctx, key)
	if err != nil {
		return nil
	}
	return student.Creator
}


// Check if the key exists in the store
func (k Keeper) StudentExists(ctx sdk.Context, key string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(types.StudentPrefix + key))
}
