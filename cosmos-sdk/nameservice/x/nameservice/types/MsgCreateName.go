package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

var _ sdk.Msg = &MsgCreateName{}

type MsgCreateName struct {
  ID      string
  Creator sdk.AccAddress `json:"creator" yaml:"creator"`
  Value string `json:"value" yaml:"value"`
  Price string `json:"price" yaml:"price"`
}

func NewMsgCreateName(creator sdk.AccAddress, value string, price string) MsgCreateName {
  return MsgCreateName{
    ID: uuid.New().String(),
		Creator: creator,
    Value: value,
    Price: price,
	}
}

func (msg MsgCreateName) Route() string {
  return RouterKey
}

func (msg MsgCreateName) Type() string {
  return "CreateName"
}

func (msg MsgCreateName) GetSigners() []sdk.AccAddress {
  return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

func (msg MsgCreateName) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg MsgCreateName) ValidateBasic() error {
  if msg.Creator.Empty() {
    return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
  }
  return nil
}