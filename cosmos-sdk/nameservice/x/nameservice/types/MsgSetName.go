package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSetName{}

type MsgSetName struct {
	Owner sdk.AccAddress `json:"owner" yaml:"owner"`
	Value string         `json:"value" yaml:"value"`
	Name  string         `json:"name" yaml:"name"`
}

func NewMsgSetName(owner sdk.AccAddress, value string, name string) MsgSetName {
	return MsgSetName{
		Owner: owner,
		Value: value,
		Name:  name,
	}
}

func (msg MsgSetName) Route() string {
	return RouterKey
}

func (msg MsgSetName) Type() string {
	return "SetName"
}

func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Owner)}
}

func (msg MsgSetName) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic 只进行无状态的校验
func (msg MsgSetName) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
	}

	if msg.Name == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name can't be empty")
	}

	if msg.Value == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "value can't be empty")
	}
	return nil
}
