<!--
order: 6
-->

# Encoding

While encoding in the SDK used to be mainly handled by `go-amino` codec, the SDK is moving towards using `gogoprotobuf` for both state and client-side encoding. {synopsis}

## Pre-requisite Readings

- [Anatomy of an SDK application](../basics/app-anatomy.md) {prereq}

## Encoding

The Cosmos SDK utilizes two binary wire encoding protocols, [Amino](https://github.com/tendermint/go-amino/) which is an object encoding specification and [Protocol Buffers](https://developers.google.com/protocol-buffers), a subset of Proto3 with an extension for
interface support. See the [Proto3 spec](https://developers.google.com/protocol-buffers/docs/proto3)
for more information on Proto3, which Amino is largely compatible with (but not with Proto2).

Due to Amino having significant performance drawbacks(缺点), being reflection-based, and
not having any meaningful cross-language/client support, Protocol Buffers, specifically
[gogoprotobuf](https://github.com/gogo/protobuf/), is being used in place of Amino.
Note, this process of using Protocol Buffers over Amino is still an ongoing process.

Binary wire encoding of types in the Cosmos SDK can be broken down into two main
categories, client encoding and store encoding. Client encoding mainly revolves
around transaction processing and signing, whereas store encoding revolves around
types used in state-machine transitions and what is ultimately stored in the Merkle
tree.

主要有两大应用: 客户端编码(client store)和存储编码(store encoding). 客户端编码主要处理交易和签名, 存储编码处理状态机中的交易,这些交易最终存储在Merkle树中.

For store encoding, protobuf definitions can exist for any type and will typically
have an Amino-based "intermediary" type. Specifically, the protobuf-based type
definition is used for serialization and persistence, whereas the Amino-based type
is used for business logic in the state-machine where they may converted back-n-forth.
Note, the Amino-based types may slowly be phased-out in the future so developers
should take note to use the protobuf message definitions where possible.


Amino编码将逐步淘汰, 所以开发者应该尽量使用 protobuf



In the `codec` package, there exists two core interfaces, `Marshaler` and `ProtoMarshaler`,
where the former encapsulates(封装) the current Amino interface except it operates on
types implementing the latter instead of generic `interface{}` types.

`Marshaler`是Amino接口; `ProtoMarshaler`是`interface{}`

In addition, there exists two implementations of `Marshaler`. The first being
`AminoCodec`, where both binary and JSON serialization is handled via Amino. The
second being `ProtoCodec`, where both binary and JSON serialization is handled
via Protobuf.


`Marshaler`有两个实现: `AmnioCodec`:通过amino处理二进制和json序列化; 和`ProtoCodec`:使用protobuf进行二进制和json的序列化.


This means that modules may use Amino or Protobuf encoding but the types must
implement `ProtoMarshaler`. If modules wish to avoid implementing this interface
for their types, they may use an Amino codec directly.

这些就是说, 模块可以使用amino或protobuf编码(需要实现ProtoMarshaler). 如果不想使用这些接口, 可以直接使用Amino


### Amino

Every module uses an Amino codec to serialize types and interfaces. This codec typically
has types and interfaces registered in that module's domain only (e.g. messages),
but there are exceptions like `x/gov`. Each module exposes a `RegisterLegacyAminoCodec` function
that allows a user to provide a codec and have all the types registered. An application
will call this method for each necessary module.

模块可以使用Amino序列化type 和 interface. 通常需要进行类型和接口的模块范围的注册. 每个模块暴露了一个`RegisterLegacyAminoCodec`函数允许用户提供codec并且拥有注册的多有类型. 应用程序会为需要的模块的调用此函数.



Where there is no protobuf-based type definition for a module (see below), Amino
is used to encode and decode raw wire bytes to the concrete type or interface:


```go
bz := keeper.cdc.MustMarshalBinaryBare(typeOrInterface)
keeper.cdc.MustUnmarshalBinaryBare(bz, &typeOrInterface)
```

Note, there are length-prefixed variants of the above functionality and this is
typically used for when the data needs to be streamed or grouped together
(e.g. `ResponseDeliverTx.Data`)





### Gogoproto

Modules are encouraged to utilize Protobuf encoding for their respective types.

推荐使用protobuf进行编码

#### FAQ

1. How to create modules using protobuf encoding?

**Defining module types**

Protobuf types can be defined to encode:
  - state
  - [`Msg`s](../building-modules/messages-and-queries.md#messages)
  - [Query services](../building-modules/query-services.md)
  - [genesis](../building-modules/genesis.md)
  
**Naming and conventions**

We encourage developers to follow industry guidelines: [Protocol Buffers style guide](https://developers.google.com/protocol-buffers/docs/style)
and [Buf](https://buf.build/docs/style-guide), see more details in [ADR 023](../architecture/adr-023-protobuf-naming.md)

2. How to update modules to protobuf encoding?

If modules do not contain any interfaces (e.g. `Account` or `Content`), then they
may simply migrate any existing types that
are encoded and persisted via their concrete Amino codec to Protobuf (see 1. for further guidelines) and accept a `Marshaler` as the codec which is implemented via the `ProtoCodec`
without any further customization.

However, if modules are to handle type interfaces, module-level .proto files should define messages which encode interfaces
using [`google.protobuf.Any`](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/any.proto).

For example, we can define `MsgSubmitEvidence` as follows where `Evidence` is
an interface:

```protobuf
// proto/cosmos/evidence/v1beta1/tx.proto

message MsgSubmitEvidence {
  string              submitter = 1;
  google.protobuf.Any evidence  = 2 [(cosmos_proto.accepts_interface) = "Evidence"];
}
```

The SDK provides support methods `MarshalAny` and `UnmarshalAny` to allow
easy encoding of state to `Any`.

Module should register interfaces using `InterfaceRegistry` which provides a mechanism for registering interfaces: `RegisterInterface(protoName string, iface interface{})` and implementations: `RegisterImplementations(iface interface{}, impls ...proto.Message)` that can be safely unpacked from Any, similarly to type registration with Amino:

+++ https://github.com/cosmos/cosmos-sdk/blob/3d969a1ffdf9a80f9ee16db9c16b8a8aa1004af6/codec/types/interface_registry.go#L23-L52

In addition, an `UnpackInterfaces` phase should be introduced to deserialization to unpack interfaces before they're needed. Protobuf types that contain a protobuf `Any` either directly or via one of their members should implement the `UnpackInterfacesMessage` interface:

```go
type UnpackInterfacesMessage interface {
  UnpackInterfaces(InterfaceUnpacker) error
}
```

#### Guidelines for protobuf message definitions

In addition to [following official guidelines](https://developers.google.com/protocol-buffers/docs/proto3#simple), we recommend to use these annotations in .proto files when dealing with interfaces:
* fields which accept interfaces should be annotated with `cosmos_proto.accepts_interface`
using the same full-qualified name passed as `protoName` to `InterfaceRegistry.RegisterInterface`
* interface implementations should be annotated with `cosmos_proto.implements_interface`
using the same full-qualified name passed as `protoName` to `InterfaceRegistry.RegisterInterface`

#### Transaction Encoding

Another important use of Protobuf is the encoding and decoding of
[transactions](./transactions.md). Transactions are defined by the application or
the SDK, but passed to the underlying consensus engine in order to be relayed to
other peers. Since the underlying consensus engine is agnostic(不可知的) to the application,
it only accepts transactions in the form of raw bytes. The encoding is done by an
object called `TxEncoder` and the decoding by an object called `TxDecoder`.

+++ https://github.com/cosmos/cosmos-sdk/blob/9ae17669d6715a84c20d52e10e2232be9f467360/types/tx_msg.go#L82-L86

```go

// TxDecoder unmarshals transaction bytes
type TxDecoder func(txBytes []byte) (Tx, error)

// TxEncoder marshals transaction to bytes
type TxEncoder func(tx Tx) ([]byte, error)
```


A standard implementation of both these objects can be found in the [`auth` module](https://github.com/cosmos/cosmos-sdk/blob/master/x/auth):

`auth`模块中的的标准实现:

+++ https://github.com/cosmos/cosmos-sdk/blob/9ae17669d6715a84c20d52e10e2232be9f467360/x/auth/tx/decoder.go

+++ https://github.com/cosmos/cosmos-sdk/blob/9ae17669d6715a84c20d52e10e2232be9f467360/x/auth/tx/encoder.go

## Next {hide}

Learn about [events](./events.md) {hide}
