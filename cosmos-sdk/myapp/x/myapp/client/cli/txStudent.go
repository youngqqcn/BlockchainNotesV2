package cli

import (
	"bufio"
    "strconv"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/username/myapp/x/myapp/types"
)

func GetCmdCreateStudent(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-student [stuname] [age] [gender] [homeaddr]",
		Short: "Creates a new student",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsStuname := string(args[0] )
			argsAge, _ := strconv.ParseInt(args[1] , 10, 64)
			argsGender, _ := strconv.ParseBool(args[2] )
			argsHomeaddr := string(args[3] )
			
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgCreateStudent(cliCtx.GetFromAddress(), string(argsStuname), int32(argsAge), bool(argsGender), string(argsHomeaddr))
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}


func GetCmdSetStudent(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-student [id]  [stuname] [age] [gender] [homeaddr]",
		Short: "Set a new student",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			argsStuname := string(args[1])
			argsAge, _ := strconv.ParseInt(args[2], 10, 64)
			argsGender, _ := strconv.ParseBool(args[3])
			argsHomeaddr := string(args[4])
			
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgSetStudent(cliCtx.GetFromAddress(), id, string(argsStuname), int32(argsAge), bool(argsGender), string(argsHomeaddr))
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdDeleteStudent(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-student [id]",
		Short: "Delete a new student by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgDeleteStudent(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
