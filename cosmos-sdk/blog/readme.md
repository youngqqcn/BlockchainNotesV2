# blog

**blog** is a blockchain application built using Cosmos SDK and Tendermint and generated with [Starport](https://github.com/tendermint/starport).

## Get started

```
starport serve
```

`serve` command installs dependencies, initializes and runs the application.

## Configure

Initialization parameters of your app are stored in `config.yml`.

### `accounts`

A list of user accounts created during genesis of your application.

| Key   | Required | Type            | Description                                       |
| ----- | -------- | --------------- | ------------------------------------------------- |
| name  | Y        | String          | Local name of the key pair                        |
| coins | Y        | List of Strings | Initial coins with denominations (e.g. "100coin") |

## Learn more

- [Starport](https://github.com/tendermint/starport)
- [Cosmos SDK documentation](https://docs.cosmos.network)
- [Cosmos Tutorials](https://tutorials.cosmos.network)
- [Channel on Discord](https://discord.gg/W8trcGV)




# 概述

`app/app.go`导入并配置了一些 sdk的模块. 并且有一个app的构造函数, app继承了 `baseapp` . app只使用了几个SDK标准的模块, 如`auth`, `bank` 和 我们自定义的`x/blog`模块

`cmd`目录下包含两个程序的源码, `blogd`全节点的程序, `blogcli`是我们的客户端程序, 主要用于查询全节点的数据和发送交易.

本示例将会使用 `key-value`存储数据, 和大多数key-value存储一样,你可以增删改查.



## 整理


- `types/TypePost.go`: 定义`Post`类型
- `types/MsgCreatePost.go`: 定义`MsgCreatePost`消息, 并实现 `sdk.Msg`消息接口
- `types/codec.go`: 为`MsgCreatePost`注册 codec
- `handlerMsgCreatePost.go`: 处理`MsgCreatePost`消息
- `handler.go`: 根据不同消息创建不同handler
- `keeper/post.go`: 用于对数据库的读和写
- `keeper/querier.go`: 查询器, 根据不同的查询请求获取数据
- `client/cli`: 用于命令行
- `client/rest`: 用户 HTTP
- `client/cli/tx.go`: 管理tx(交易)相关的子命令
- `client/cli/txPost.go`: 管理 post相关交易的子命令(如: create-post),发布post需要发起一笔新的交易
- `client/cli/query.go`: 管理查询相关的子命令
- `client/cli/queryPost.go`:管理从数据库中查询post的命令 (如: list-post )

