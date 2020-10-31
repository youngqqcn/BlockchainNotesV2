<!--
order: 7
-->

# Events

`Event`s are objects that contain information about the execution of the application. They are mainly used by service providers like block explorers and wallet to track the execution of various messages and index transactions. {synopsis}

`Event`包含了应用程序执行信息. 主要用于服务提供者,像区块浏览器和钱包追踪不同的message执行和交易的索引. 


> 提问: 提供是否可以提供websocket? 以msg为基础, 实现具体的具体的应用场景. 代替智能合约?

## Pre-requisite Readings

- [Anatomy of an SDK application](../basics/app-anatomy.md) {prereq}

## Events

Events are implemented in the Cosmos SDK as an alias of the ABCI `Event` type and
take the form of: `{eventType}.{eventAttribute}={value}`.

+++ https://github.com/tendermint/tendermint/blob/bc572217c07b90ad9cee851f193aaa8e9557cbc7/abci/types/types.pb.go#L2187-L2193

```go
type Event struct {
	Type                 string          `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Attributes           []common.KVPair `protobuf:"bytes,2,rep,name=attributes,proto3" json:"attributes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}
```


Events contain:

- A `type`, which is meant to categorize an event at a high-level (e.g. by module or action).
- A list of `attributes`, which are key-value pairs that give more information about
  the categorized `event`.
  +++ https://github.com/cosmos/cosmos-sdk/blob/7d7821b9af132b0f6131640195326aa02b6751db/types/events.go#L51-L56
  
  ```go
    // Attribute defines an attribute wrapper where the key and value are
    // strings instead of raw bytes.
	Attribute struct {
		Key   string `json:"key"`
		Value string `json:"value,omitempty"`
	}

  ```

Events are returned to the underlying consensus engine in the response of the following ABCI messages:

- [`BeginBlock`](./baseapp.md#beginblock)
- [`EndBlock`](./baseapp.md#endblock)
- [`CheckTx`](./baseapp.md#checktx)
- [`DeliverTx`](./baseapp.md#delivertx)

Events, the `type` and `attributes`, are defined on a **per-module basis** in the module's
`/types/events.go` file, and triggered from the module's [`handler`](../building-modules/handler.md)
via the [`EventManager`](#eventmanager). In addition, each module documents its events under
`spec/xx_events.md`.

## EventManager

In Cosmos SDK applications, events are managed by an abstraction called the `EventManager`.
Internally, the `EventManager` tracks a list of `Events` for the entire execution flow of a
transaction or `BeginBlock`/`EndBlock`.

+++ https://github.com/cosmos/cosmos-sdk/blob/7d7821b9af132b0f6131640195326aa02b6751db/types/events.go#L16-L20

```go
// EventManager implements a simple wrapper around a slice of Event objects that
// can be emitted from.
type EventManager struct {
	events Events
}
```


The `EventManager` comes with a set of useful methods to manage events. Among them, the one that is
used the most by module and application developers is the `EmitEvent` method, which tracks
an `event` in the `EventManager`.

`EventManager`有一些有用的方法来管理event. 其中用的最多的是`EmitEvent`方法, 它可以用来跟踪`EvnetManager`的`event`



+++ https://github.com/cosmos/cosmos-sdk/blob/7d7821b9af132b0f6131640195326aa02b6751db/types/events.go#L29-L31


```go
// EmitEvent stores a single Event object.
func (em *EventManager) EmitEvent(event Event) {
	em.events = em.events.AppendEvent(event)
}
```

Module developers should handle event emission via the `EventManager#EmitEvent` in each message
`Handler` and in each `BeginBlock`/`EndBlock` handler. The `EventManager` is accessed via
the [`Context`](./context.md), where event emission generally follows this pattern:

```go
ctx.EventManager().EmitEvent(
    sdk.NewEvent(eventType, sdk.NewAttribute(attributeKey, attributeValue)),
)
```

Module's `handler` function should also set a new `EventManager` to the `context` to isolate emitted events per `message`:

模块的handler函数应该为`context`设置一个新的`EventManager`以是每个message分开(isolate)触发event


```go
func NewHandler(keeper Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
        ctx = ctx.WithEventManager(sdk.NewEventManager())
        switch msg := msg.(type) {
```

See the [`Handler`](../building-modules/handler.md) concept doc for a more detailed
view on how to typically implement `Events` and use the `EventManager` in modules.

## Subscribing to Events

It is possible to subscribe to `Events` via Tendermint's [Websocket](https://tendermint.com/docs/app-dev/subscribing-to-events-via-websocket.html#subscribing-to-events-via-websocket).
This is done by calling the `subscribe` RPC method via Websocket:

订阅事件: 可以通过Tendermint的websocket订阅事件

```json
{
    "jsonrpc": "2.0",
    "method": "subscribe",
    "id": "0",
    "params": {
        "query": "tm.event='eventCategory' AND eventType.eventAttribute='attributeValue'"
    }
}
```

The main `eventCategory` you can subscribe to are:

- `NewBlock`: Contains `events` triggered during `BeginBlock` and `EndBlock`.
- `Tx`: Contains `events` triggered during `DeliverTx` (i.e. transaction processing). 
- `ValidatorSetUpdates`: Contains validator set updates for the block.

> 思考: 可以指定eventType和attributeValue, 那么可不可以通过websocket订阅某个账户的交易事件?


These events are triggered from the `state` package after a block is committed. You can get the
full list of `event` categories [here](https://godoc.org/github.com/tendermint/tendermint/types#pkg-constants).

这些event在区块提交后由`state`包触发. 

```go

const (
    // Block level events for mass consumption by users.
    // These events are triggered from the state package,
    // after a block has been committed.
    // These are also used by the tx indexer for async indexing.
    // All of this data can be fetched through the rpc.
    EventNewBlock            = "NewBlock"
    EventNewBlockHeader      = "NewBlockHeader"
    EventNewEvidence         = "NewEvidence"
    EventTx                  = "Tx"
    EventValidatorSetUpdates = "ValidatorSetUpdates"

    // Internal consensus events.
    // These are used for testing the consensus state machine.
    // They can also be used to build real-time consensus visualizers.
    EventCompleteProposal = "CompleteProposal"
    EventLock             = "Lock"
    EventNewRound         = "NewRound"
    EventNewRoundStep     = "NewRoundStep"
    EventPolka            = "Polka"
    EventRelock           = "Relock"
    EventTimeoutPropose   = "TimeoutPropose"
    EventTimeoutWait      = "TimeoutWait"
    EventUnlock           = "Unlock"
    EventValidBlock       = "ValidBlock"
    EventVote             = "Vote"
)

```


The `type` and `attribute` value of the `query` allow you to filter the specific `event` you are looking for. For example, a `transfer` transaction triggers an `event` of type `Transfer` and has `Recipient` and `Sender` as `attributes` (as defined in the [`events` file of the `bank` module](https://github.com/cosmos/cosmos-sdk/blob/master/x/bank/types/events.go)). Subscribing to this `event` would be done like so:

```go
// bank module event types
const (
	EventTypeTransfer = "transfer"

	AttributeKeyRecipient = "recipient"
	AttributeKeySender    = "sender"

	AttributeValueCategory = ModuleName
)
```


可以通过事件过滤转账事件

```json
{
    "jsonrpc": "2.0",
    "method": "subscribe",
    "id": "0",
    "params": {
        "query": "tm.event='Tx' AND transfer.sender='senderAddress'"
    }
}
```

where `senderAddress` is an address following the [`AccAddress`](../basics/accounts.md#addresses) format.

## Next {hide}

Learn about SDK [telemetry](./telemetry.md) {hide}
