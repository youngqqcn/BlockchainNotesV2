package scavenge

import (
	"crypto/sha256"
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"

	"github.com/youngqqcn/scavenge/x/scavenge/keeper"
	"github.com/youngqqcn/scavenge/x/scavenge/types"
)

func handleMsgRevealSolution(ctx sdk.Context, k keeper.Keeper, msg types.MsgRevealSolution) (*sdk.Result, error) {

	solutionScavengerBytes := []byte(msg.Solution + msg.Scavenger.String())
	solutionScavengerHash := sha256.Sum256(solutionScavengerBytes)
	solutionScavengerHashString := hex.EncodeToString(solutionScavengerHash[:])

	_, err := k.GetCommit(ctx, solutionScavengerHashString)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Commit with that hash doesn't exists")
	}

	solutionHash := sha256.Sum256([]byte(msg.Solution))
	solutionHashString := hex.EncodeToString(solutionHash[:])

	// 从数据库中获取此 scavenge  的信息
	scavenge, err := k.GetScavenge(ctx, solutionHashString)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Scavenge with that solution hash doesn't exists")
	}

	// 判断 scavenge 是否已经被解决, 如果已经被解决了, 报错
	if scavenge.Scavenger != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Scavenge has already been solved")
	}
	
	// 如果此题尚未被解决, 则设置解题人, 答案
	scavenge.Scavenger = msg.Scavenger
	scavenge.Solution = msg.Solution

	moduleAcc := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
	sdkError := k.CoinKeeper.SendCoins(ctx, moduleAcc, scavenge.Scavenger, scavenge.Reward)
	if sdkError != nil {
		return nil, sdkError
	}

	k.SetScavenge(ctx, scavenge)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeSolveScavenge),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Scavenger.String()),
			sdk.NewAttribute(types.AttributeSolutionHash, solutionHashString),
			sdk.NewAttribute(types.AttributeDescription, scavenge.Description),
			sdk.NewAttribute(types.AttributeSolution, msg.Solution),
			sdk.NewAttribute(types.AttributeScavenger, scavenge.Scavenger.String()),
			sdk.NewAttribute(types.AttributeReward, scavenge.Reward.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
