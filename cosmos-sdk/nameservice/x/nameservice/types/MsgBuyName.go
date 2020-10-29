package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgBuyName{}

type MsgBuyName struct {
	Buyer sdk.AccAddress `json:"creator" yaml:"creator"`
	Bid   sdk.Coins      `json:"bid" yaml:"bid"`
	Name  string         `json:"name" yaml:"name"`
}

func NewMsgBuyName(Buyer sdk.AccAddress, name string, bid sdk.Coins) MsgBuyName {
	return MsgBuyName{
		Buyer: Buyer,
		Bid:   bid,
		Name:  name,
	}
}

func (msg MsgBuyName) Route() string {
	return RouterKey
}

func (msg MsgBuyName) Type() string {
	return "BuyName"
}

func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Buyer)}
}

func (msg MsgBuyName) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic  只能进行无状态检查
func (msg MsgBuyName) ValidateBasic() error {
	if msg.Buyer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "buyer can't be empty")
	}

	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}

	if !msg.Bid.IsAllPositive() {
		return sdkerrors.ErrInsufficientFunds
	}

	return nil
}
