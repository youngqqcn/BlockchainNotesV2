package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgDeleteName{}

type MsgDeleteName struct {
	Owner sdk.AccAddress `json:"owner" yaml:"owner"`
	Name  string         `json:"name" ymal:"name"`
}

func NewMsgDeleteName(owner sdk.AccAddress, name string) MsgDeleteName {
	return MsgDeleteName{
		Owner: owner,
		Name:  name,
	}
}

func (msg MsgDeleteName) Route() string {
	return RouterKey
}

func (msg MsgDeleteName) Type() string {
	return "DeleteName"
}

func (msg MsgDeleteName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Owner)}
}

func (msg MsgDeleteName) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

//ValidateBasic 注意: 只能进行无状态检查, 不能去操作数据库
func (msg MsgDeleteName) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Owner can't be empty")
	}

	if msg.Name == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name can't be empty")
	}

	return nil
}
