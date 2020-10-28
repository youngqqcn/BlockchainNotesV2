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

func GetCmdCommitSolution(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-commit [solution] ",
		Short: "Creates a solution for scavenge",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSolution := string(args[0])
			// argsSolutionScavengerHash := string(args[1])
			solutionHash := sha256.Sum256([]byte(argsSolution))
			solutionHashHex := hex.EncodeToString(solutionHash[:])

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var scavenger = cliCtx.GetFromAddress().String()
			solutionScavengerHash := sha256.Sum256([]byte(argsSolution + scavenger))
			solutionScavengerHashHex := hex.EncodeToString(solutionScavengerHash[:])

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgCommitSolution(cliCtx.GetFromAddress(), solutionHashHex, solutionScavengerHashHex)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
