package myapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/username/myapp/x/myapp/types"
	"github.com/username/myapp/x/myapp/keeper"
)

func handleMsgCreateStudent(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateStudent) (*sdk.Result, error) {
	var student = types.Student{
		Creator: msg.Creator,
		ID:      msg.ID,
    	Stuname: msg.Stuname,
    	Age: msg.Age,
    	Gender: msg.Gender,
    	Homeaddr: msg.Homeaddr,
	}
	k.CreateStudent(ctx, student)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
