package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

// 用于检查 MsgCreatePost 是否实现了sdk.Msg接口的方法
var _ sdk.Msg = &MsgCreatePost{}

type MsgCreatePost struct {
	ID      string
	Creator sdk.AccAddress `json:"creator" yaml:"creator"`
	Title   string         `json:"title" yaml:"title"`
	Body    string         `json:"body" yaml:"body"`
}

// 构造函数
func NewMsgCreatePost(creator sdk.AccAddress, title string, body string) MsgCreatePost {
	return MsgCreatePost{
		ID:      uuid.New().String(),
		Creator: creator,
		Title:   title,
		Body:    body,
	}
}

// 以下函数是实现sdk.Msg接口
// type Msg interface {
// 	Route() string
// 	Type() string
// 	ValidateBasic() error
// 	GetSignBytes() []byte
// 	GetSigners() []AccAddress
// }

// Return the message type.
// Must be alphanumeric or empty.
func (msg MsgCreatePost) Route() string {
	return RouterKey
}

// Returns a human-readable string for the message, intended for utilization
// within tags
func (msg MsgCreatePost) Type() string {
	return "CreatePost"
}

// Signers returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (msg MsgCreatePost) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Get the canonical byte representation of the Msg.
func (msg MsgCreatePost) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic does a simple validation check that
// doesn't require access to any other information.
func (msg MsgCreatePost) ValidateBasic() error {
	if msg.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
	}
	return nil
}
