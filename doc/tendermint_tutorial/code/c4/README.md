C4
# 非对称加密与身份识别

了解采用非对称加密技术实现去中心化身份识别的机制，掌握密钥与地址的
生成方法，以及身份验证的实现。

## 运行预置代码

### 1、secp256k1密钥与地址生成

在2#终端执行如下命令：

```
~/repo/go/src/hubwiz.com/c4$ go run key-secp256k1.go
```

### 2、ed25519密钥与地址生成

在2#终端执行如下命令：

```
~/repo/go/src/hubwiz.com/c4$ go run key-ed25519.go
```

### 3、数据签名与验证

在2#终端执行如下命令：

```
~/repo/go/src/hubwiz.com/c4$ go run sign-verify.go
```

