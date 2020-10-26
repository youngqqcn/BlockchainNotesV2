package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgDeleteStudent{}

type MsgDeleteStudent struct {
  ID      string         `json:"id" yaml:"id"`
  Creator sdk.AccAddress `json:"creator" yaml:"creator"`
}

func NewMsgDeleteStudent(id string, creator sdk.AccAddress) MsgDeleteStudent {
  return MsgDeleteStudent{
    ID: id,
		Creator: creator,
	}
}

func (msg MsgDeleteStudent) Route() string {
  return RouterKey
}

func (msg MsgDeleteStudent) Type() string {
  return "DeleteStudent"
}

func (msg MsgDeleteStudent) GetSigners() []sdk.AccAddress {
  return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

func (msg MsgDeleteStudent) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg MsgDeleteStudent) ValidateBasic() error {
  if msg.Creator.Empty() {
    return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
  }
  return nil
}