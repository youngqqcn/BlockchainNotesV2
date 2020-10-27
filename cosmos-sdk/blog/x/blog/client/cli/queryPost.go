package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	blogtypes "github.com/youngqqcn/blog/x/blog/types"
)

func GetCmdListPost(queryRoute string, cdc *codec.Codec) *cobra.Command {

	return &cobra.Command{
		Use:   "list-post",
		Short: "list all post",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/"+blogtypes.QueryListPost, queryRoute), nil)
			if err != nil {
				fmt.Printf("list Post error: %s\n", err.Error())
			}

			var output []blogtypes.Post
			cdc.MustUnmarshalJSON(res, &output)
			return ctx.PrintOutput(output)
		},
	}

}
