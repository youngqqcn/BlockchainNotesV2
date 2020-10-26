package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSetStudent{}

type MsgSetStudent struct {
  ID      string      `json:"id" yaml:"id"`
  Creator sdk.AccAddress `json:"creator" yaml:"creator"`
  Stuname string `json:"stuname" yaml:"stuname"`
  Age int32 `json:"age" yaml:"age"`
  Gender bool `json:"gender" yaml:"gender"`
  Homeaddr string `json:"homeaddr" yaml:"homeaddr"`
}

func NewMsgSetStudent(creator sdk.AccAddress, id string, stuname string, age int32, gender bool, homeaddr string) MsgSetStudent {
  return MsgSetStudent{
    ID: id,
		Creator: creator,
    Stuname: stuname,
    Age: age,
    Gender: gender,
    Homeaddr: homeaddr,
	}
}

func (msg MsgSetStudent) Route() string {
  return RouterKey
}

func (msg MsgSetStudent) Type() string {
  return "SetStudent"
}

func (msg MsgSetStudent) GetSigners() []sdk.AccAddress {
  return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

func (msg MsgSetStudent) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg MsgSetStudent) ValidateBasic() error {
  if msg.Creator.Empty() {
    return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
  }
  return nil
}