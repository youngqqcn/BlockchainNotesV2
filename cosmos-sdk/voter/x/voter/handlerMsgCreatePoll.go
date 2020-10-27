package voter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/youngqqcn/voter/x/voter/keeper"
	"github.com/youngqqcn/voter/x/voter/types"
)

func handleMsgCreatePoll(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreatePoll) (*sdk.Result, error) {

	// 新建一轮投票, 需要支付币, 接收方是临时创建的一个地址
	moduleAddress := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

	amount, err := sdk.ParseCoins("200yqq")
	if err != nil {
		return nil, err
	}

	err = k.CoinKeeper.SendCoins(ctx, msg.Creator, moduleAddress, amount)
	if err != nil {
		return nil, err
	}

	var poll = types.Poll{
		Creator: msg.Creator,
		ID:      msg.ID,
		Title:   msg.Title,
		Options: msg.Options,
	}
	k.CreatePoll(ctx, poll)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
