package types

// nameservice module event types
const (
	// TODO: Create your event types
	// EventType<Action>    		= "action"
	EventTypeSetName    = "SetName"
	EventTypeBuyName    = "BuyName"
	EventTypeDeleteName = "DeleteName"

	// TODO: Create keys fo your events, the values will be derivided from the msg
	// AttributeKeyAddress  		= "address"
	AttributeName  = "name"
	AttributeValue = "value"
	AttributePrice = "price"
	AttributeOwner = "owner"
	AttributeBuyer= "buyer"
	AttributeBid = "bid"

	// TODO: Some events may not have values for that reason you want to emit that something happened.
	// AttributeValueDoubleSign = "double_sign"

	AttributeValueCategory = ModuleName
)
