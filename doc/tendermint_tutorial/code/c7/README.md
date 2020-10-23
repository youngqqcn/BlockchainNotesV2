# 代币案例：使用多版本状态库

了解tendermint多版本数据库的实现机制与用途，学习在abci应用中如何使用
多版本数据库。


目录文件组织：

- daemon.go: abci应用
- cli.go: 节点客户端
- lib: 公用代码目录
- iavl-demo.go：iavl测试代码
- store-demo.go：状态库封装测试
- wallet：钱包文件
- account.db：iavl库目录

## 预置代码运行

### 1、多版本状态库测试

在2#终端首先执行以下命令清除原有的库目录：

```
~/repo/go/src/hubwiz.com/c7$ rm -rf account.db
```

然后在2#终端运行测试代码：

```
~/repo/go/src/hubwiz.com/c7$ go run iavl-demo.go
```

### 2、iavl封装代码测试

在2#终端运行测试代码：

```
~/repo/go/src/hubwiz.com/c7$ go run store-demo.go
```

### 3、ABCI应用

在2#终端启动ABCI应用：

```
~/repo/go/src/hubwiz.com/c7$ go run daemon.go
```

在1#终端重新初始化并启动tendermint

```
~$ tendermint unsafe_reset_all
~$ tendermint node
```

在3#终端执行客户端程序的子命令，例如：

发行代币：

```
~/repo/go/src/hubwiz.com/c7$ go run cli.go issue-tx
```

转账：

```
~/repo/go/src/hubwiz.com/c7$ go run cli.go transfer-tx
```

查询账户michael的余额：

```
~/repo/go/src/hubwiz.com/c7$ go run cli.go query michael
```

  