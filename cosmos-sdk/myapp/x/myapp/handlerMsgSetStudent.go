package myapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/username/myapp/x/myapp/types"
	"github.com/username/myapp/x/myapp/keeper"
)

func handleMsgSetStudent(ctx sdk.Context, k keeper.Keeper, msg types.MsgSetStudent) (*sdk.Result, error) {
	var student = types.Student{
		Creator: msg.Creator,
		ID:      msg.ID,
    	Stuname: msg.Stuname,
    	Age: msg.Age,
    	Gender: msg.Gender,
    	Homeaddr: msg.Homeaddr,
	}
	if !msg.Creator.Equals(k.GetStudentOwner(ctx, msg.ID)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}

	k.SetStudent(ctx, student)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
