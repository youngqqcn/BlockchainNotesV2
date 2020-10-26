package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

var _ sdk.Msg = &MsgCreateStudent{}

type MsgCreateStudent struct {
  ID      string
  Creator sdk.AccAddress `json:"creator" yaml:"creator"`
  Stuname string `json:"stuname" yaml:"stuname"`
  Age int32 `json:"age" yaml:"age"`
  Gender bool `json:"gender" yaml:"gender"`
  Homeaddr string `json:"homeaddr" yaml:"homeaddr"`
}

func NewMsgCreateStudent(creator sdk.AccAddress, stuname string, age int32, gender bool, homeaddr string) MsgCreateStudent {
  return MsgCreateStudent{
    ID: uuid.New().String(),
		Creator: creator,
    Stuname: stuname,
    Age: age,
    Gender: gender,
    Homeaddr: homeaddr,
	}
}

func (msg MsgCreateStudent) Route() string {
  return RouterKey
}

func (msg MsgCreateStudent) Type() string {
  return "CreateStudent"
}

func (msg MsgCreateStudent) GetSigners() []sdk.AccAddress {
  return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

func (msg MsgCreateStudent) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg MsgCreateStudent) ValidateBasic() error {
  if msg.Creator.Empty() {
    return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
  }
  return nil
}