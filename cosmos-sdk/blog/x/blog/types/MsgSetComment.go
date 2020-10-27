package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSetComment{}

type MsgSetComment struct {
  ID      string      `json:"id" yaml:"id"`
  Creator sdk.AccAddress `json:"creator" yaml:"creator"`
  Body string `json:"body" yaml:"body"`
  PostID string `json:"postID" yaml:"postID"`
}

func NewMsgSetComment(creator sdk.AccAddress, id string, body string, postID string) MsgSetComment {
  return MsgSetComment{
    ID: id,
		Creator: creator,
    Body: body,
    PostID: postID,
	}
}

func (msg MsgSetComment) Route() string {
  return RouterKey
}

func (msg MsgSetComment) Type() string {
  return "SetComment"
}

func (msg MsgSetComment) GetSigners() []sdk.AccAddress {
  return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

func (msg MsgSetComment) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg MsgSetComment) ValidateBasic() error {
  if msg.Creator.Empty() {
    return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
  }
  return nil
}