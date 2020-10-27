package voter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/youngqqcn/voter/x/voter/keeper"
	"github.com/youngqqcn/voter/x/voter/types"
)

func handleMsgCreateVote(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateVote) (*sdk.Result, error) {

	// // 投票需要支付币, 接收方是临时创建的一个地址
	// moduleAddress := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

	// amount, err := sdk.ParseCoins("200yqq")
	// if err != nil {
	// 	return nil, err
	// }

	// err = k.CoinKeeper.SendCoins(ctx, msg.Creator, moduleAddress, amount)
	// if err != nil {
	// 	return nil, err
	// }

	var vote = types.Vote{
		Creator: msg.Creator,
		ID:      msg.ID,
		PollID:  msg.PollID,
		Value:   msg.Value,
	}
	k.CreateVote(ctx, vote)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
