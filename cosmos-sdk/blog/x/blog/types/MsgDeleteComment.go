package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgDeleteComment{}

type MsgDeleteComment struct {
  ID      string         `json:"id" yaml:"id"`
  Creator sdk.AccAddress `json:"creator" yaml:"creator"`
}

func NewMsgDeleteComment(id string, creator sdk.AccAddress) MsgDeleteComment {
  return MsgDeleteComment{
    ID: id,
		Creator: creator,
	}
}

func (msg MsgDeleteComment) Route() string {
  return RouterKey
}

func (msg MsgDeleteComment) Type() string {
  return "DeleteComment"
}

func (msg MsgDeleteComment) GetSigners() []sdk.AccAddress {
  return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

func (msg MsgDeleteComment) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg MsgDeleteComment) ValidateBasic() error {
  if msg.Creator.Empty() {
    return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
  }
  return nil
}