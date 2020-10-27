package pofe

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/youngqqcn/pofe/x/pofe/keeper"
	"github.com/youngqqcn/pofe/x/pofe/types"
)

func handleMsgCreateClaim(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateClaim) (*sdk.Result, error) {
	var claim = types.Claim{
		Creator: msg.Creator,
		// ID:      msg.ID,
		Proof: msg.Proof,
	}
	k.CreateClaim(ctx, claim)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
