package scavenge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"

	"github.com/youngqqcn/scavenge/x/scavenge/keeper"
	"github.com/youngqqcn/scavenge/x/scavenge/types"
)

func handleMsgCreateScavenge(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateScavenge) (*sdk.Result, error) {
	var scavenge = types.Scavenge{
		Creator: msg.Creator,
		// ID:      msg.ID,
		Description:  msg.Description,
		SolutionHash: msg.SolutionHash,
		Reward:       msg.Reward,
		// Solution:     msg.,
		// Scavenger: msg.Scavenger,
	}

	_, err := k.GetScavenge(ctx, scavenge.SolutionHash)
	if err == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Scavenge with that solution hash already exists")
	}

	moduleAcc := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
	sdkError := k.CoinKeeper.SendCoins(ctx, scavenge.Creator, moduleAcc, scavenge.Reward)
	if sdkError != nil {
		return nil, sdkError
	}

	k.CreateScavenge(ctx, scavenge)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeCreateScavenge),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
			sdk.NewAttribute(types.AttributeDescription, msg.Description),
			sdk.NewAttribute(types.AttributeSolutionHash, msg.SolutionHash),
			sdk.NewAttribute(types.AttributeReward, msg.Reward.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
