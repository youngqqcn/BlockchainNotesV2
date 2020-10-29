package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/youngqqcn/nameservice/x/nameservice/keeper"
	"github.com/youngqqcn/nameservice/x/nameservice/types"
	// abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k keeper.Keeper /* TODO: Define what keepers the module needs */, data types.GenesisState) {
	// TODO: Define logic for when you would like to initalize a new genesis

	// 从genesis状态加载到keeper中
	for _, record := range data.WhoisRecords {
		k.SetWhois(ctx, record.Value, record)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) (data types.GenesisState) {
	// TODO: Define logic for exporting state

	var records []types.Whois

	it := k.GetNamesIterator(ctx)
	for ; it.Valid(); it.Next() {
		name := string(it.Key())
		whois, _ := k.GetWhois(ctx, name)
		records = append(records, whois)
	}
	return types.NewGenesisState(records)
}
