

>参考文档:
https://github.com/tendermint/starport/tree/develop/docs

安装 starport

- 直接使用npm安装: `npm i -g @tendermint/starport` 
- 或从源码编译: `git clone https://github.com/tendermint/starport && cd starport && make `



```

starport app github.com/username/myapp && cd myapp

starport serve

starport type post title body

```

启动结果:

```
yqq@ubuntu:myapp$ starport serve
Cosmos' version is: Launchpad

📦 Installing dependencies...
🛠️  Building the app...
🙂 Created an account. Password (mnemonic): duck secret spy velvet entire knock venue forward boring ability pulp alcohol wrestle lecture disease sorry host whip picture home address song year hockey //第一个账户的助记词
🙂 Created an account. Password (mnemonic): craft learn nose habit future panda faculty beef harsh butter deny share aunt measure whip hazard damage horn include theory negative animal music profit  //第二个账户的助记词
🌍 Running a Cosmos 'myapp' app with Tendermint at http://0.0.0.0:26657.
🌍 Running a server at http://0.0.0.0:1317 (LCD)

🚀 Get started: http://localhost:12345

```

浏览器打开 :
-  http://localhost:12345/#/ 可以看到一个类似区块链浏览器的页面
- http://localhost:8080/  可以进行转账(需要进行登录, 使用助记词进行登录,启动时有输出), 还可以发布文章(title post)


### 配置:

注意: 如果修改`config.yml`中的`chainid`, 也要修改 `~/.myappcli/config/config.toml` 中的`chain_id` , 否则客户端转账会出错

如下, 可以修改代币的名称

```
version: 1
accounts:
  - name: user1
    coins: ["1000yqq", "100000000stake"]
  - name: user2
    coins: ["500yqq"]
validator:
  name: user1
  staked: "100000000stake"

genesis:
  chain_id: "yqqchain"
  app_state:
    staking:
      params:
        bond_denom: "stake"
```


### 地址

可以修改 `app/prefix.go` 中的 `AccountAddressPrefix` 修改地址的前缀


### 进行转账

```
#查询账户列表
myappcli keys list    

#转账 10yqq
myappcli tx send yqq1emcl5gwm7mxx0agc0ut47n697na3p7v6fhgul6 yqq1qlm38vzszgxsmarudt2kghdlm3huc9435c62nw 10yqq

#查询账户余额
myappcli query account yqq1qlm38vzszgxsmarudt2kghdlm3huc9435c62nw

#查询交易
yqq@ubuntu:myapp$ myappcli query tx 4250BCE8F1B24A801A94AFC194B2231E695BC3CB0F4D3E02729247C9CA05214D
{
  "height": "30",
  "txhash": "4250BCE8F1B24A801A94AFC194B2231E695BC3CB0F4D3E02729247C9CA05214D",
  "raw_log": "[{\"msg_index\":0,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"yqq1xypzmjw5mwcntgs3wgg0l34v667rr05hjhem4g\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"yqq1h0nxg67w9ruzmtackw5euhrtdswrzy64a8kjw5\"},{\"key\":\"sender\",\"value\":\"yqq1xypzmjw5mwcntgs3wgg0l34v667rr05hjhem4g\"},{\"key\":\"amount\",\"value\":\"10yqq\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "send"
            },
            {
              "key": "sender",
              "value": "yqq1xypzmjw5mwcntgs3wgg0l34v667rr05hjhem4g"
            },
            {
              "key": "module",
              "value": "bank"
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "recipient",
              "value": "yqq1h0nxg67w9ruzmtackw5euhrtdswrzy64a8kjw5"
            },
            {
              "key": "sender",
              "value": "yqq1xypzmjw5mwcntgs3wgg0l34v667rr05hjhem4g"
            },
            {
              "key": "amount",
              "value": "10yqq"
            }
          ]
        }
      ]
    }
  ],
  "gas_wanted": "200000",
  "gas_used": "47231",
  "tx": {
    "type": "cosmos-sdk/StdTx",
    "value": {
      "msg": [
        {
          "type": "cosmos-sdk/MsgSend",
          "value": {
            "from_address": "yqq1xypzmjw5mwcntgs3wgg0l34v667rr05hjhem4g",
            "to_address": "yqq1h0nxg67w9ruzmtackw5euhrtdswrzy64a8kjw5",
            "amount": [
              {
                "denom": "yqq",
                "amount": "10"
              }
            ]
          }
        }
      ],
      "fee": {
        "amount": [],
        "gas": "200000"
      },
      "signatures": [
        {
          "pub_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A9+o1hpdHjP+dKu+QbOvudFMmdbnGXwM7QZ3kHIn5RUs"
          },
          "signature": "zySV68VTyKlKx0LhmW7k1B5W/mVipOuVMl2/RWyVW05+mpXS3Qk5uLplBeamt86IHMGA9VljQgY1xQ88RYoKFg=="
        }
      ],
      "memo": ""
    }
  },
  "timestamp": "2020-10-26T07:35:25Z"
}




#发布文章
myappcli tx  myapp create-post  thisisposttest thisiscontent  --from yqq1xypzmjw5mwcntgs3wgg0l34v667rr05hjhem4g 
{
  "chain_id": "yqqchain",
  "account_number": "2",
  "sequence": "2",
  "fee": {
    "amount": [],
    "gas": "200000"
  },
  "msgs": [
    {
      "type": "myapp/CreatePost",
      "value": {
        "ID": "49f97cad-26ca-45f9-a3da-0ee93cebdf11",
        "creator": "yqq1xypzmjw5mwcntgs3wgg0l34v667rr05hjhem4g",
        "title": "thisisposttest",
        "body": "thisiscontent"
      }
    }
  ],
  "memo": ""
}

confirm transaction before signing and broadcasting [y/N]: y
{
  "height": "0",
  "txhash": "1F5529B6BBDA3B4CB26D243F3EAAAF3DF087F8B08F281ECAA034F60D0AFD76E2",
  "raw_log": "[]"
}



#查询文章
yqq@ubuntu:myapp$ myappcli query myapp get-post 49f97cad-26ca-45f9-a3da-0ee93cebdf11
{
  "creator": "yqq1xypzmjw5mwcntgs3wgg0l34v667rr05hjhem4g",
  "id": "49f97cad-26ca-45f9-a3da-0ee93cebdf11",
  "title": "thisisposttest",
  "body": "thisiscontent"
}


```



## 架构

### 介绍

Serve

在项目目录中运行 `starport serve` 命令会安装依赖,构建并初始化应用程序.然后启动Tendermint RPC(`localhost:26657`) , 还有 LCD(local client daemon) `localhost:1317`




Key-Value 存储

使用 `starport type` 可以自动生成相应的代码(`messages`, `handlers`, `keepers`, `CLI`和 `REST`, 以及 相应的数据类型定义)

例如: `starport type post title body` 就是创建了一个类型`Post`, 包含两个字段: `title` 和 `body`



实际操作: 创建一个 `student` type :

```

starport type student stuname  age:int gender:bool  homeaddr

重启服务


#创建 student
yqq@ubuntu:cosmos-sdk$ myappcli tx  myapp create-student  yqq 10 true china --from user1
{
  "chain_id": "yqqchain",
  "account_number": "2",
  "sequence": "1",
  "fee": {
    "amount": [],
    "gas": "200000"
  },
  "msgs": [
    {
      "type": "myapp/CreateStudent",
      "value": {
        "ID": "ce6ec2cd-0648-4049-b4ae-8340526db077",
        "creator": "yqq1hz44g5n4ufjzlh5wt8jnpq3a9rhaamm6wxx34u",
        "stuname": "yqq",
        "age": 10,
        "gender": true,
        "homeaddr": "china"
      }
    }
  ],
  "memo": ""
}

confirm transaction before signing and broadcasting [y/N]: y
{
  "height": "0",
  "txhash": "46072D1BFFC92525CF599222C4608568AC9FCC8A3346659C4BD48EE4AF0BC313",
  "raw_log": "[]"
}



#查询交易
$myappcli query tx 46072D1BFFC92525CF599222C4608568AC9FCC8A3346659C4BD48EE4AF0BC313

{
....
"tx": {
    "type": "cosmos-sdk/StdTx",
    "value": {
      "msg": [
        {
          "type": "myapp/CreateStudent",
          "value": {
            "ID": "ce6ec2cd-0648-4049-b4ae-8340526db077",
            "creator": "yqq1hz44g5n4ufjzlh5wt8jnpq3a9rhaamm6wxx34u",
            "stuname": "yqq",
            "age": 10,
            "gender": true,
            "homeaddr": "china"
          }
        }
      ],
      "fee": {
        "amount": [],
        "gas": "200000"
      },
      "signatures": [
        {
          "pub_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "Ai9oUTZD1IM1YJGRaSjL16Ngi/zSQdumGTfz3BSVXJMS"
          },
          "signature": "fkSLAnGzjSqy9DDHhLhnSOS5DenVjFAJB8aWFO1m7ul9uvbHzI5+u7B5j8D+Wbgqiq9lbJ6RXt1cED1SJDuh0Q=="
        }
      ],
      "memo": ""
    }
....
}

# 获取student
yqq@ubuntu:cosmos-sdk$ myappcli query myapp list-student
[
  {
    "creator": "yqq1hz44g5n4ufjzlh5wt8jnpq3a9rhaamm6wxx34u",
    "id": "ce6ec2cd-0648-4049-b4ae-8340526db077",
    "stuname": "yqq",
    "age": 10,
    "gender": true,
    "homeaddr": "china"
  }
]

yqq@ubuntu:cosmos-sdk$ myappcli query myapp get-student ce6ec2cd-0648-4049-b4ae-8340526db077
{
  "creator": "yqq1hz44g5n4ufjzlh5wt8jnpq3a9rhaamm6wxx34u",
  "id": "ce6ec2cd-0648-4049-b4ae-8340526db077",
  "stuname": "yqq",
  "age": 10,
  "gender": true,
  "homeaddr": "china"   
}

```

### 目录结构

- `abci.go`: 模块的 `BeginBlocker` and `EndBlocker` 的实现(如果有的话).
- `client/`: 模块 `CLI` and `REST`客户端实现和测试.
- `exported/`:  定义了基于接口依赖的"合约", 来保持依赖关系的清晰干净
- `handler.go`: 模块的消息处理handler
- `keeper/`: 模块的keeper实现,和一些辅助功能的实现
- `types/`: 模块的类型定义, 例如:消息, KVStore,参数类型,Protocal Buffer的定义,和 expected_keepers.go合约
- `module.go`: 模块的 `AppModule` and `AppModuleBasic`接口实现




### 标准模块

#### Auth

`auth`模块负责区块链上的账户,以及基础的交易类型.CosmosSDK默认提供了交易,手续费,签名和重放攻击保护, 还有账户授权.  账户授权和staking通常是分开使用的


#### Bank

`bank`模块负责代币转账, 已经验证转账人的有效性. 还负责检查代币的总发行量(所有账户余额的总和)



#### Staking

`staking`模块可以构建高级的 PoS(Proof of Stake)系统. 可以创建验证者(Validators)也可以向验证者进行委托



##### Distribution

`distribution`模块负责代币的增发.当一个新的代币创建时, 将分发给验证者和委托人, 验证者会收取一些佣金.当创建验证者时创建这可以设置佣金, 这个佣金可以修改.



#### Params

`params`模块主要负责保存一些全局的参数, 这些参数会在节点运行时变化, 可以通过通过`government`模块进行修改, 当多数股东都同意时, 参数的修改才会生效.




以上这些模块, 使用starport创建项目时会自动添加,  除了这些模块之外, 还有一些其他模块

#### 使用模块


可以使用    `starport module create modulename` 添加模块. 当手动添加一个模块之后, 需要修改 `app/app.go` 和 `myappcli/main.go` . starport 可以方便进行模块增加和编辑 .

`starport module import <modulename>` 可以导入命令




### 编写自定义模块

使用 starport 可以很方便添加自己的模块 , 前面已经使用过的命令 `starport type xxxxx` , 这个命令自动帮我们创建了 `handler`, `types`, 和 `message`

如果没有 starport , 也可以手动添加模块, 理解创建模块的步骤和原理是必须.

#### 类型

- 创建的基础类型在`types/typeXXXXX.go`, 例如student 模块类型: 

    ```go
    type Student struct {
        Creator sdk.AccAddress `json:"creator" yaml:"creator"`
        ID      string         `json:"id" yaml:"id"`
        Stuname string `json:"stuname" yaml:"stuname"`
        Age int32 `json:"age" yaml:"age"`
        Gender bool `json:"gender" yaml:"gender"`
        Homeaddr string `json:"homeaddr" yaml:"homeaddr"`
    }

    ```


- 消息类型定义在`types/MsgCreateXXXXX` , 例如student模块的消息类型:

    ```go
    ...

    type MsgCreateStudent struct {
    ID      string
    Creator sdk.AccAddress `json:"creator" yaml:"creator"`
    Stuname string `json:"stuname" yaml:"stuname"`
    Age int32 `json:"age" yaml:"age"`
    Gender bool `json:"gender" yaml:"gender"`
    Homeaddr string `json:"homeaddr" yaml:"homeaddr"`
    }
    ...

    ```

- `client/rest`: 目录定义了暴露的一些 REST API

- `client/cli`: 目录定义了一些命令行交互的功能


#### 前端

starport 创建的项目, 自带了一个默认的前端, 使用 `vue`框架编写. 可以以此为基础构建你自己的前端




# 教程


[Poll](https://github.com/cosmos/sdk-tutorials/blob/master/voter/index.md)

[Blog](https://github.com/cosmos/sdk-tutorials/blob/master/blog/tutorial/01-index.md)

[Scavenge](https://github.com/cosmos/sdk-tutorials/blob/master/scavenge/tutorial/01-background.md)

[Nameservice](https://github.com/cosmos/sdk-tutorials/blob/master/nameservice/tutorial/00-intro.md)







