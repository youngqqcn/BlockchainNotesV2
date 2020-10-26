package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSetPost{}

type MsgSetPost struct {
  ID      string      `json:"id" yaml:"id"`
  Creator sdk.AccAddress `json:"creator" yaml:"creator"`
  Title string `json:"title" yaml:"title"`
  Body string `json:"body" yaml:"body"`
}

func NewMsgSetPost(creator sdk.AccAddress, id string, title string, body string) MsgSetPost {
  return MsgSetPost{
    ID: id,
		Creator: creator,
    Title: title,
    Body: body,
	}
}

func (msg MsgSetPost) Route() string {
  return RouterKey
}

func (msg MsgSetPost) Type() string {
  return "SetPost"
}

func (msg MsgSetPost) GetSigners() []sdk.AccAddress {
  return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

func (msg MsgSetPost) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg MsgSetPost) ValidateBasic() error {
  if msg.Creator.Empty() {
    return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
  }
  return nil
}