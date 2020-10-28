package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgRevealSolution struct {
	Scavenger    sdk.AccAddress `json:"scavenger" yaml:"scavenger"`
	SolutionHash string         `json:"solutionHash" yaml:"solutionHash"`
	Solution     string         `json:"solution" yaml:"solution"`
}

func NewMsgRevealSolution(scavenger sdk.AccAddress, solution string) MsgRevealSolution {
	var solutionHash = sha256.Sum256([]byte(solution))
	var solutionHashString = hex.EncodeToString(solutionHash[:])

	return MsgRevealSolution{
		Scavenger:    scavenger,
		SolutionHash: solutionHashString,
		Solution:     solution,
	}
}

const RevealSolutionConst = "RevealSolution"

func (msg MsgRevealSolution) Route() string {
	return RouterKey
}

func (msg MsgRevealSolution) Type() string {
	return RevealSolutionConst
}
func (msg MsgRevealSolution) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Scavenger)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgRevealSolution) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic 为 AnteHandler 提供基础的校验
func (msg MsgRevealSolution) ValidateBasic() error {

	if msg.Scavenger.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "scavenger can't be empty")
	}

	if msg.SolutionHash == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "solutionScavengerHash can't be empty")
	}

	if msg.Solution == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "solution can't be empty")
	}

	// 对结果的进行校验

	solutionHash := sha256.Sum256([]byte(msg.Solution))
	solutionHashString := hex.EncodeToString(solutionHash[:])

	if msg.SolutionHash != solutionHashString {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress,
			fmt.Sprintf("hash or solution(%s) do not equal solutionHash(%s)\n",
				msg.SolutionHash, solutionHashString))
	}
	return nil
}
