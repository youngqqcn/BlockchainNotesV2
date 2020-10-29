package nameservice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/youngqqcn/nameservice/x/nameservice/keeper"
	"github.com/youngqqcn/nameservice/x/nameservice/types"
)

func handleMsgBuyName(ctx sdk.Context, k keeper.Keeper, msg types.MsgBuyName) (*sdk.Result, error) {

	price, _ := k.GetPrice(ctx, msg.Name)
	// if err != nil {
		// 域名尚未卖出, 尚未存在价格
		// return nil
	// }

	// 如果报价小于当前的域名的价格, 则购买失败
	if  msg.Bid.IsAllLTE( price ) { 
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "bid must be greater than current bid ")
	}

	if k.HasOwner(ctx, msg.Name) {

		owner := k.GetOwner(ctx, msg.Name)

		// 如果购买人和owner相同, 
		if msg.Buyer.Equals( owner  ) {
			return nil, fmt.Errorf("buyer is owner of name")
		}

		// 从买方的账户 转给 域名所有者
		err := k.CoinKeeper.SendCoins(ctx, msg.Buyer, owner , msg.Bid)
		if err != nil {
			return nil, err
		}
	} else {
		// 直接从买方账户扣除
		_, err :=  k.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid)
		if err != nil {
			return nil, err
		}
	}

	whois := types.NewWhois()
	whois.Price = msg.Bid
	whois.Owner = msg.Buyer

	k.SetWhois(ctx, msg.Name, whois)

	ctx.EventManager().Events().AppendEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeSetName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer.String()),
			sdk.NewAttribute(types.AttributeBuyer, msg.Buyer.String()),
			sdk.NewAttribute(types.AttributeName, msg.Name),
			sdk.NewAttribute(types.AttributeBid, msg.Bid.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
