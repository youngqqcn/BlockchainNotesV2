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


下面我们开始创建一个类似博客的应用程序, 第一步, 手动创建一个  `Post` 类型:

示例代码: https://github.com/cosmos/sdk-tutorials/blob/master/blog/blog/x/blog/types/TypePost.go



