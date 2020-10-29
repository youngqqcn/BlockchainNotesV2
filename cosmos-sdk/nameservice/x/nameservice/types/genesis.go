package types

import "fmt"

// GenesisState - all nameservice state that must be provided at genesis
type GenesisState struct {
	// TODO: Fill out what is needed by the module for genesis
	WhoisRecords []Whois `json:"whoisRecords"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(records []Whois /* TODO: Fill out with what is needed for genesis state */) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		WhoisRecords: nil,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state, these values will be initialized but empty
		WhoisRecords: []Whois{},
	}
}

// ValidateGenesis validates the nameservice genesis parameters
func ValidateGenesis(data GenesisState) error {
	// TODO: Create a sanity check to make sure the state conforms to the modules needs
	for _, record := range data.WhoisRecords {

		// 对数据进行基本的检查,正常的话不会手动导入genesis 状态是由上次运行保存的
		if record.Owner == nil {
			return fmt.Errorf("invalid WhoisRecords: Value: %s, missing Owner ", record.Value)
		}
		if len(record.Value) == 0 {
			return fmt.Errorf("invalid WhoisRecords: Owner: %s, missing Value", record.Owner)
		}
		if record.Price == nil {
			return fmt.Errorf("invalid WhoisRecords: Value: %s, missing Price", record.Price)
		}
	}

	return nil
}
