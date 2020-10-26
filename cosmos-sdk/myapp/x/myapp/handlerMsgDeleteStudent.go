package myapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/username/myapp/x/myapp/types"
	"github.com/username/myapp/x/myapp/keeper"
)

// Handle a message to delete name
func handleMsgDeleteStudent(ctx sdk.Context, k keeper.Keeper, msg types.MsgDeleteStudent) (*sdk.Result, error) {
	if !k.StudentExists(ctx, msg.ID) {
		// replace with ErrKeyNotFound for 0.39+
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msg.ID)
	}
	if !msg.Creator.Equals(k.GetStudentOwner(ctx, msg.ID)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	k.DeleteStudent(ctx, msg.ID)
	return &sdk.Result{}, nil
}
