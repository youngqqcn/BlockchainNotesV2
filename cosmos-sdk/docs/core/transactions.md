<!--
order: 2
-->

# Transactions

`Transactions` are objects created by end-users to trigger state changes in the application. {synopsis}

`Transactions`是一些由用户创建的对象, 可以触发应用程序的状态改变.

## Pre-requisite Readings

* [Anatomy of an SDK Application](../basics/app-anatomy.md) {prereq}

## Transactions

Transactions are comprised of metadata held in [contexts](./context.md) and [messages](../building-modules/messages-and-queries.md) that trigger state changes within a module through the module's [Handler](../building-modules/handler.md). 

`Transactions`由contexts中的元数据(metadata)和messages组成.它可以通过模块的handler触发状态(state)改变.

When users want to interact with an application and make state changes (e.g. sending coins), they create transactions. Each of a transaction's `message`s must be signed using the private key associated with the appropriate account(s), before the transaction is broadcasted to the network. A transaction must then be included in a block, validated, and approved by the network through the consensus process. To read more about the lifecycle of a transaction, click [here](../basics/tx-lifecycle.md).

当用户和应用程序交互时并且做一些状态改变(比如:代币转账), 他们会创建交易. 每个交易的`message`s 必须在广播前使用账户的私钥签名. 一个交易被打包进区块,验证,  通过共识被区块链网络所执行. 


## Type Definition

Transaction objects are SDK types that implement the `Tx` interface

`Transaction`类型是 SDK 类型, 实现了 `Tx`接口

+++ https://github.com/cosmos/cosmos-sdk/blob/7d7821b9af132b0f6131640195326aa02b6751db/types/tx_msg.go#L34-L41


```go
// Transactions objects must fulfill the Tx
type Tx interface {
    // Gets the all the transaction's messages.
    // 获取transactions的messages, 一个transaction可以包含多个message
	GetMsgs() []Msg

	// ValidateBasic does a simple and lightweight validation check that doesn't
    // require access to any other information.
    // 例如 auth 模块的 StdTx ValidateBasic函数会检查交易的签名是否正确,以及手续费是否超额
    // 要注意和 message的ValidateBasic函数区分. message的只是检查message.
    // 当 baseapp.runTx 执行时, 首先会使用 auth  使用每个message的ValidateBasic进行检查, 然后才会
    // 使用transaction的ValidateBasic进行检查
	ValidateBasic() Error
}

// TxDecoder unmarshals transaction bytes
type TxDecoder func(txBytes []byte) (Tx, Error)

// TxEncoder marshals transaction to bytes
type TxEncoder func(tx Tx) ([]byte, error)

```

It contains the following methods:

* **GetMsgs:** unwraps the transaction and returns a list of its message(s) - one transaction may have one or multiple [messages](../building-modules/messages-and-queries.md#messages), which are defined by module developers.
* **ValidateBasic:** includes lightweight, [*stateless*](../basics/tx-lifecycle.md#types-of-checks) checks used by ABCI messages [`CheckTx`](./baseapp.md#checktx) and [`DeliverTx`](./baseapp.md#delivertx) to make sure transactions are not invalid. For example, the [`auth`](https://github.com/cosmos/cosmos-sdk/tree/master/x/auth) module's `StdTx` `ValidateBasic` function checks that its transactions are signed by the correct number of signers and that the fees do not exceed what the user's maximum. Note that this function is to be distinct from the `ValidateBasic` functions for *`messages`*, which perform basic validity checks on messages only. For example, when [`runTx`](./baseapp.md#runtx) is checking a transaction created from the [`auth`](https://github.com/cosmos/cosmos-sdk/tree/master/x/auth/spec) module, it first runs `ValidateBasic` on each message, then runs the `auth` module AnteHandler which calls `ValidateBasic` for the transaction itself.
* **TxEncoder:** Nodes running the consensus engine (e.g. Tendermint Core) are responsible for gossiping(分发) transactions and ordering them into blocks, but only handle them in the generic `[]byte` form. Transactions are always [marshaled](./encoding.md) (encoded) before they are relayed to nodes, which compacts them to facilitate gossiping and helps maintain the consensus engine's separation from from application logic. The Cosmos SDK allows developers to specify any deterministic encoding format for their applications; the default is Amino.
* **TxDecoder:** [ABCI](https://tendermint.com/docs/spec/abci/) calls from the consensus engine to the application, such as `CheckTx` and `DeliverTx`, are used to process transaction data to determine validity and state changes. Since transactions are passed in as `txBytes []byte`, they need to first be unmarshaled (decoded) using `TxDecoder` before any logic is applied.

The most used implementation of the `Tx` interface is  [`StdTx` from the `auth` module](https://github.com/cosmos/cosmos-sdk/blob/master/x/auth/types/stdtx.go). As a developer, using `StdTx` as your transaction format is as simple as importing the `auth` module in your application (which can be done in the [constructor of the application](../basics/app-anatomy.md#constructor-function)) 


用的最多的`Tx`接口的实现时`auth`模块的`StdTx`类型.开发者只需要在你的应用程序中导入`auth`模块即可使用, 可以在app的构造函数中完成.


## Transaction Process

A transaction is created by an end-user through one of the possible [interfaces](#interfaces). In the process, two contexts and an array of [messages](#messages) are created, which are then used to [generate](#transaction-generation) the transaction itself. The actual state changes triggered by transactions are enabled by the [handlers](#handlers). The rest of the document will describe each of these components, in this order.

### CLI and REST Interfaces

Application developers create entrypoints to the application by creating a [command-line interface](../interfaces/cli.md) and/or [REST interface](../interfaces/rest.md), typically found in the application's `./cmd` folder. These interfaces allow users to interact with the application through command-line or through HTTP requests.

For the [command-line interface](../building-modules/module-interfaces.md#cli), module developers create subcommands to add as children to the application top-level transaction command `TxCmd`. For [HTTP requests](../building-modules/module-interfaces.md#legacy-rest), module developers specify acceptable request types, register REST routes, and create HTTP Request Handlers.

When users interact with the application's interfaces, they invoke(调用) the underlying(底层的) modules' handlers or command functions, directly creating messages.

### Messages

**`Message`s** are module-specific objects that trigger state transitions within the scope of the module they belong to. Module developers define the `message`s for their module by implementing the `Msg` interface, and also define a [`Handler`](../building-modules/handler.md) to process them.


`Message`s 是模块特定的对象, 它可以在所属的模块范围内触发状态转变. 模块的开发者为他们的模块定义`message`并实现`Msg`接口, 还要定义一个`Handler`处理这个消息.



+++ https://github.com/cosmos/cosmos-sdk/blob/7d7821b9af132b0f6131640195326aa02b6751db/types/tx_msg.go#L8-L29


```go

// Transactions messages must fulfill the Msg
type Msg interface {

	// Return the message type.
    // Must be alphanumeric or empty.
    // 返回消息的类型, 必须时字母和数字组成,  
    // 主要用于消息路由
	Route() string

	// Returns a human-readable string for the message, intended for utilization
	// within tags
	Type() string

	// ValidateBasic does a simple validation check that
    // doesn't require access to any other information.
    // 提供基本检查
	ValidateBasic() Error

    // Get the canonical byte representation of the Msg.
    // 用于签名
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
    // CONTRACT: Returns addrs in some deterministic order.
    // 所有签名必须有效
    // 签名者必须以确定的顺序返回
	GetSigners() []AccAddress
}
```



`Message`s in a module are typically defined in a `msgs.go` file (though not always), and one handler with multiple functions to handle each of the module's `message`s is defined in a `handler.go` file.

Note: module `messages` are not to be confused with [ABCI Messages](https://tendermint.com/docs/spec/abci/abci.html#messages) which define interactions between the Tendermint and application layers.


注意: 不要把模块的`message`和 ABCI Message 搞混了, ABCI定义了Tendermint和应用层的交互.



To learn more about `message`s, click [here](../building-modules/messages-and-queries.md#messages).

While messages contain the information for state transition logic, a transaction's other metadata and relevant information are stored in the `TxBuilder` and `Context`.

### Transaction Generation

Transactions are first created by end-users through an `appcli tx` command through the command-line or a POST request to an HTTPS server. For details about transaction creation, click [here](../basics/tx-lifecycle.md#transaction-creation).

[`Contexts`](https://godoc.org/context) are immutable objects that contain all the information needed to process a request. In the process of creating a transaction through the `auth` module (though it is not mandatory to create transactions this way), two contexts are created: the [`Context`](../interfaces/query-lifecycle.md#context) and `TxBuilder`. Both are automatically generated and do not need to be defined by application developers, but do require input from the transaction creator (e.g. using flags through the CLI).

`Contexts`一个不可变对象, 它包含了处理一个请求所需要的全部的信息. 在通过`auth`模块创建交易的过程中, 两个context会被创建: `Context`和`TxBuilder`. 两个都是自动生成的, 并且不需要应用程序开发者定义, 但是需要输入交易的创建者.(例如.使用命令需要加 --from=xxx)

The `TxBuilder` contains data closely related with the processing of transactions.

+++ https://github.com/cosmos/cosmos-sdk/blob/7d7821b9af132b0f6131640195326aa02b6751db/x/auth/types/txbuilder.go#L18-L31

```go

// TxBuilder implements a transaction context created in SDK modules.
type TxBuilder struct {
	txEncoder          sdk.TxEncoder
	keybase            crkeys.Keybase
	accountNumber      uint64
	sequence           uint64
	gas                uint64
	gasAdjustment      float64
	simulateAndExecute bool
	chainID            string
	memo               string
	fees               sdk.Coins
	gasPrices          sdk.DecCoins
}
```


* `TxEncoder` defined by the developer for this type of transaction. Used to encode messages before being processed by nodes running Tendermint.
* `Keybase` that manages the user's keys and is used to perform signing operations.
* `AccountNumber` from which this transaction originated.
* `Sequence`, the number of transactions that the user has sent out, used to prevent replay attacks(重放攻击). 
* `Gas` option chosen by the users for how to calculate how much gas they will need to pay. A common option is "auto" which generates an automatic estimate.
* `GasAdjustment` to adjust the estimate of gas by a scalar value, used to avoid underestimating the amount of gas required.
* `SimulateAndExecute` option to simply simulate the transaction execution without broadcasting.
* `ChainID` representing which blockchain this transaction pertains to.
* `Memo` to send with the transaction.
* `Fees`, the maximum amount the user is willing to pay in fees. Alternative to specifying gas prices.
* `GasPrices`, the amount per unit of gas the user is willing to pay in fees. Alternative to specifying fees.

The `Context` is initialized using the application's `codec` and data more closely related to the user interaction with the interface, holding data such as the output to the user and the broadcast mode. Read more about `Context` [here](../interfaces/query-lifecycle.md#context).

Every message in a transaction must be signed by the addresses specified by `GetSigners`. The signing process must be handled by a module, and the most widely used one is the [`auth`](https://github.com/cosmos/cosmos-sdk/tree/master/x/auth/spec) module. Signing is automatically performed when the transaction is created, unless the user choses to generate and sign separately. The `TxBuilder` (namely, the `KeyBase`) is used to perform the signing operations, and the `Context` is used to broadcast transactions.

`TxBuilder`用来执行交易签名的操作, `Context`用来广播交易.


`auth/client/utils/tx.go`命令行源码: 

```go
// CompleteAndBroadcastTxCLI implements a utility function that facilitates
// sending a series of messages in a signed transaction given a TxBuilder and a
// QueryContext. It ensures that the account exists, has a proper number and
// sequence set. In addition, it builds and signs a transaction with the
// supplied messages. Finally, it broadcasts the signed transaction to a node.
func CompleteAndBroadcastTxCLI(txBldr authtypes.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) error {
	txBldr, err := PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	fromName := cliCtx.GetFromName()

	// 模拟交易
	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err = EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return err
		}

		gasEst := GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	// 是否模拟交易, 如果模拟交易(看看gas是否超出)
	if cliCtx.Simulate {
		return nil
	}

	// 进行确认操作
	if !cliCtx.SkipConfirm {
		stdSignMsg, err := txBldr.BuildSignMsg(msgs)
		if err != nil {
			return err
		}

		var json []byte
		if viper.GetBool(flags.FlagIndentResponse) {
			json, err = cliCtx.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
			if err != nil {
				panic(err)
			}
		} else {
			json = cliCtx.Codec.MustMarshalJSON(stdSignMsg)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return err
		}
	}

	// 构建交易并签名
	// build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(fromName, keys.DefaultKeyPass, msgs)
	if err != nil {
		return err
	}

	// 广播交易
	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	return cliCtx.PrintOutput(res)
}

```



### Handlers

Since `message`s are module-specific types, each module needs a [`handler`](../building-modules/handler.md) to process all of its `message` types and trigger state changes within the module's scope. This design puts more responsibility on module developers, allowing application developers to reuse common functionalities without having to implement state transition logic repetitively. To read more about `handler`s, click [here](../building-modules/handler.md).

因为`message`是特定模块的类型, 每个模块需要一个`handler`处理它的所有的`message`类型, 并且在模块范围内触发状态更改. 这样设计, 将更多的责任给了模块的开发者, 允许应用程序开发者重用公共的功能而不需要重复实现状态转换的逻辑 .

## Next {hide}

Learn about the [context](./context.md) {hide}
