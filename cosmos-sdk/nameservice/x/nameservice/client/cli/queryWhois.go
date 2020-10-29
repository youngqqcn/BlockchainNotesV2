package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/youngqqcn/nameservice/x/nameservice/types"
)

func GetCmdListWhois(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list-whois",
		Short: "list all whois",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/"+types.QueryListWhois, queryRoute), nil)
			if err != nil {
				fmt.Printf("could not list Whois\n%s\n", err.Error())
				return nil
			}

			// var out []types.Whois
			var out []types.QueryWhois
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdGetWhois(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-whois [key]",
		Short: "Query a whois by key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			key := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryGetWhois, key), nil)
			if err != nil {
				fmt.Printf("could not found whois: %s , error:%s\n", key, err.Error())
				return nil
			}

			var out types.Whois
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdResolveName(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve-name [name]",
		Short: "resolve a name ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := string(args[0])

			if len(name) == 0 {
				fmt.Printf("could not resolve empty name\n")
				return nil
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryResolveName, name), nil)
			if err != nil {
				fmt.Printf("could not resolve name %s \n%s\n", name, err.Error())
				return nil
			}

			var out types.Whois
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}