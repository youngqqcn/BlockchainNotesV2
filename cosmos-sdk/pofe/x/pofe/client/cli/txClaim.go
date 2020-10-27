package cli

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/os"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/youngqqcn/pofe/x/pofe/types"
)

func GetCmdCreateClaim(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-claim [path-to-file]",
		Short: "Creates a new claim from a path to a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			// 获取文件路径
			filePath := string(args[0])
			if !os.FileExists(filePath) {
				return fmt.Errorf("%s not exists", filePath)
			}

			// 读取文件内容
			bz, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}

			// 对文件内容进行sha256
			hasher := sha256.New()
			hasher.Write(bz)

			// 将hash编码为十六进制字符串
			proofHashHex := hex.EncodeToString(hasher.Sum(nil))

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgCreateClaim(cliCtx.GetFromAddress(), string(proofHashHex))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdSetClaim(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-claim [id]  [proof]",
		Short: "Set a new claim",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			argsProof := string(args[1])

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgSetClaim(cliCtx.GetFromAddress(), id, string(argsProof))
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdDeleteClaim(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-claim [id]",
		Short: "Delete a new claim by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgDeleteClaim(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
