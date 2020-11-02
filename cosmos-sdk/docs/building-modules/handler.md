<!--
order: 4
-->

# Handlers

A `Handler` designates(指定) a function that processes [`message`s](./messages-and-queries.md#messages). `Handler`s are specific to the module in which they are defined, and only process `message`s defined within the said module. They are called from `baseapp` during [`DeliverTx`](../core/baseapp.md#delivertx). {synopsis}

## Pre-requisite Readings

- [Module Manager](./module-manager.md) {prereq}
- [Messages and Queries](./messages-and-queries.md) {prereq}

## `handler` type

The `handler` type defined in the Cosmos SDK specifies the typical structure of a `handler` function.

+++ https://github.com/cosmos/cosmos-sdk/blob/7d7821b9af132b0f6131640195326aa02b6751db/types/handler.go#L4

```go
// Handler defines the core of the state transition function of an application.
type Handler func(ctx Context, msg Msg) Result
```

Let us break it down:

- The [`Msg`](./messages-and-queries.md#messages) is the actual object being processed. 
- The [`Context`](../core/context.md) contains all the necessary information needed to process the `msg`, as well as a cache-wrapped copy of the latest state. If the `msg` is succesfully processed, the modified version of the temporary state contained in the `ctx` will be written to the main state.
- The [`*Result`] returned to `baseapp`, which contains (among other things) information on the execution of the `handler` and [`events`](../core/events.md).
	+++ https://github.com/cosmos/cosmos-sdk/blob/d55c1a26657a0af937fa2273b38dcfa1bb3cff9f/proto/cosmos/base/abci/v1beta1/abci.proto#L81-L95

## Implementation of a module `handler`s

Module `handler`s are typically implemented in a `./handler.go` file inside the module's folder. The
[module manager](./module-manager.md) is used to add the module's `handler`s to the
[application's `router`](../core/baseapp.md#message-routing) via the `Route()` method. Typically,
the manager's `Route()` method simply constructs a Route that calls a `NewHandler()` method defined in `handler.go`,
which looks like the following:

```go
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *MsgType1:
			return handleMsgType1(ctx, keeper, msg)

		case *MsgType2:
			return handleMsgType2(ctx, keeper, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}
```

First, the `handler` function sets a new `EventManager` to the context to isolate events per `msg`.
Then, this simple switch returns a `handler` function specific to the type of the received `message`. These `handler` functions are the ones that actually process `message`s, and usually follow the following 2 steps:

- First, they perform *stateful* checks to make sure the `message` is valid. At this stage, the `message`'s `ValidateBasic()` method has already been called, meaning *stateless* checks on the message (like making sure parameters are correctly formatted) have already been performed. Checks performed in the `handler` can be more expensive and require access to the state. For example, a `handler` for a `transfer` message might check that the sending account has enough funds to actually perform the transfer. To access the state, the `handler` needs to call the [`keeper`'s](./keeper.md) getter functions. 
- Then, if the checks are successfull, the `handler` calls the [`keeper`'s](./keeper.md) setter functions to actually perform the state transition. 


`handler`处理`message`主要有2步:

- 第一: 执行 *有状态* 检查, 以确保`message`是有效的. 到这个阶段, `message`的`ValidateBasic`早已经被调用并执行了*无状态*检查(主要检查参数格式是否正确). 执行有状态检查时很昂贵的,因为需要访问状态(数据库).例如, 处理`transfer`消息需要检查发送方的余额是否足够. 为了访问状态, `handler`需要调用`keeper`的getter函数
- 第二: 如果检查成功, 则`handler`会调用 `keeper`的setter方法, 执行状态修改操作.


这里一投票的demo为例:

```go

func handleMsgSetPoll(ctx sdk.Context, k keeper.Keeper, msg types.MsgSetPoll) (*sdk.Result, error) {
	var poll = types.Poll{
		Creator: msg.Creator,
		ID:      msg.ID,
    	Title: msg.Title,
    	Options: msg.Options,
    }
    
    // 执行有状态的检查
	if !msg.Creator.Equals(k.GetPollOwner(ctx, msg.ID)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}

    // 通过keeper改变状态
	k.SetPoll(ctx, poll)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

```



Before returning, `handler` functions generally emit one or multiple [`events`](../core/events.md) via the `EventManager` held in the `ctx`:

```go
ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,  // e.g. sdk.EventTypeMessage for a message, types.CustomEventType for a custom event defined in the module
			sdk.NewAttribute(attributeKey, attributeValue),
		),
    )
```

These `events` are relayed back to the underlying consensus engine and can be used by service providers to implement services around the application. Click [here](../core/events.md) to learn more about `events`. 

Finally, the `handler` function returns a `*sdk.Result` which contains the aforementioned `events` and an optional `Data` field. 

+++ https://github.com/cosmos/cosmos-sdk/blob/d55c1a26657a0af937fa2273b38dcfa1bb3cff9f/proto/cosmos/base/abci/v1beta1/abci.proto#L81-L95

Next is an example of how to return a `*Result` from the `gov` module:

+++ https://github.com/cosmos/cosmos-sdk/blob/d55c1a26657a0af937fa2273b38dcfa1bb3cff9f/x/gov/handler.go#L67-L70

For a deeper look at `handler`s, see this [example implementation of a `handler` function](https://github.com/cosmos/cosmos-sdk/blob/d55c1a26657a0af937fa2273b38dcfa1bb3cff9f/x/gov/handler.go) from the `gov` module.

The `handler` can then be registered from [`AppModule.Route()`](./module-manager.md#appmodule) as shown in the example below:

+++ https://github.com/cosmos/cosmos-sdk/blob/228728cce2af8d494c8b4e996d011492139b04ab/x/gov/module.go#L143-L146

## Telemetry

New [telemetry metrics](../core/telemetry.md) can be created from the `handler` when handling messages for instance. 

This is an example from the `auth` module:

+++ https://github.com/cosmos/cosmos-sdk/blob/d55c1a26657a0af937fa2273b38dcfa1bb3cff9f/x/auth/vesting/handler.go#L68-L80

## Next {hide}

Learn about [query services](./query-services.md) {hide}
