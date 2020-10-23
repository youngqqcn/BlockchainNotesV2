# 案例：发行自己的代币

使用tendermint实现一个代币应用。

目录文件组织：

- daemon.go: abci应用
- cli.go: 节点客户端
- lib: 公用代码目录
- token-sm.go：状态机测试代码
- token-wallet：钱包测试代码
- token-tx.go：交易结构测试代码
- token-codec.go：编解码测试
- wallet：钱包文件

## 使用预置代码


### 1、账户状态机测试

在2#终端执行以下命令：

```
~/repo/go/src/hubwiz.com/c5$ go run token-sm.go
```

### 2、钱包测试

在2#终端执行如下命令初始化钱包：

```
~/repo/go/src/hubwiz.com/c5$ go run token-wallet.go init
```

在2#终端执行如下命令载入钱包：

```
~/repo/go/src/hubwiz.com/c5$ go run token-wallet.go load
```

### 3、交易结构测试

在2#终端执行如下命令：

```
~/repo/go/src/hubwiz.com/c5$ go run token-tx.go
```

### 4、编解码器测试

在2#终端执行如下命令：

```
~/repo/go/src/hubwiz.com/c5$ go run token-codec.go
```

### 5、ABCI应用

在2#终端启动ABCI应用：

```
~/repo/go/src/hubwiz.com/c5$ go run daemon.go
```

在1#终端重新初始化并启动tendermint

```
~$ tendermint unsafe_reset_all
~$ tendermint node
```

在3#终端执行客户端程序的子命令，例如：

发行代币：

```
~/repo/go/src/hubwiz.com/c5$ go run cli.go issue-tx
```

转账：

```
~/repo/go/src/hubwiz.com/c5$ go run cli.go transfer-tx
```

查询账户michael的余额：

```
~/repo/go/src/hubwiz.com/c5$ go run cli.go query michael
```


