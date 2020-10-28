package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSetName{}

type MsgSetName struct {
  ID      string      `json:"id" yaml:"id"`
  Creator sdk.AccAddress `json:"creator" yaml:"creator"`
  Value string `json:"value" yaml:"value"`
  Price string `json:"price" yaml:"price"`
}

func NewMsgSetName(creator sdk.AccAddress, id string, value string, price string) MsgSetName {
  return MsgSetName{
    ID: id,
		Creator: creator,
    Value: value,
    Price: price,
	}
}

func (msg MsgSetName) Route() string {
  return RouterKey
}

func (msg MsgSetName) Type() string {
  return "SetName"
}

func (msg MsgSetName) GetSigners() []sdk.AccAddress {
  return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

func (msg MsgSetName) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg MsgSetName) ValidateBasic() error {
  if msg.Creator.Empty() {
    return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
  }
  return nil
}