package cli

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/youngqqcn/scavenge/x/scavenge/types"
)

func GetCmdCreateScavenge(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-scavenge [description] [solution] [reward] ",
		Short: "Creates a new scavenge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsDescription := string(args[0])
			argsSolution := string(args[1])
			argsReward := string(args[2])

			solutionHash := sha256.Sum256([]byte(argsSolution))
			solutionHashHex := hex.EncodeToString(solutionHash[:])

			reward, err := sdk.ParseCoins(argsReward)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgCreateScavenge(cliCtx.GetFromAddress(), string(argsDescription), string(solutionHashHex), reward, string(argsSolution))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdRevealSolution(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "reveal-solution [solution]",
		Short: "reveal a solution for scavenge",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			var solution = args[0]
			msg := types.NewMsgRevealSolution(cliCtx.GetFromAddress(), solution)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}

}
