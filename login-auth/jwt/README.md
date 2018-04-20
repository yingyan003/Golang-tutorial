# go实现jwt

# 基础知识前瞻
------------

> 前言：
因为http协议是一种无状态协议，这就涉及用户访问系统的状态保持问题。

在登录系统的设计中，当用户访问系统时，服务器需要认证用户登录相关的信息，以决定用户能否登录到系统中。
这一设计一般有2种实现方式：session和jwt。每种开发语言都有其相应的实现包，本文主要介绍jwt在go中的实践。

## jwt简介

JSON Web Token(jwt)是一种规范，常用于用户与服务器间的认证。

### jwt结构

jwt由以下三部分构成：
* Header:头部 （对应：Header）
* Claims:声明  (对应：Payload)
* Signature:签名  (对应：Signature)

#### Header
Header中指明jwt的签名算法，如
```
{
  "typ": "JWT",
  "alg": "HS256"
}
```

#### Claims
声明中有jwt自身预置的，使用时可选。当然，我们也可以加入自定义的声明，
如uid，userName之类信息，但一定不要声明重要或私密的信息，因为这些信息是可破解的。

#### Signature
在生成jwt的token（令牌的意思）串时，先将Header和Claims用base64编码,再用Header中指定的加密算法，
将编码后的2个字符串进行加密（签名）。加密时需要用到一个signString签名串，我们可指定自己的signString，
不同的signString生成的加密结果不一样（解密时可能也需要同样的串，视加密算法而定）。

最后生成的jwt token串格式是：Header.Claims.Signature.如

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.
TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ
```
> （实际的token串不换行，这里只为展示清晰）

## session VS jwt

* 二者最主要的区别是用户登录状态的保存机制。

session：将用户登录的状态信息保存在服务器，客户端只存储uid。
客户端访问服务器时，服务器拿着uid去获取相应的服务器session(一个数据结构)，并判断用户登录状态是否有效，是则允许登录系统，否则返回客户端（重新登录）。

jwt：将用户登录状态信息保存在客户端（即token串，因为token只保存在客户端），token中可以设置token失效时间。
客户端访问服务器时，带着该token，服务器解析token，解析成功且登录状态有效可以放行，否则返回客户端（具体逻辑视个人情况而定）。

* 对服务器压力

session：因为session会话信息保存在服务器，会增加服务器I/O压力

jwt token：因服务器需要解析token（如base64解码，解密），会增加服务器计算压力。
但token中可以保存一些用户的基本信息，服务器解析即可获取，免去了查数据库的必要

2种方式各有各的好，至于使用哪种方式，视自身情况而定。

## 扩展
> 本文重点在于jwt的go实现，概念只做大致介绍

对于初次接触登录系统设计或开发的同学，先对相关概念和设计思路有个大致的了解，对于快速熟悉和开发系统有很大的帮助。

可以依次参考以下连接：

* jwt了解
http://blog.leapoahead.com/2015/09/06/understanding-jwt/
* 八幅漫画理解jwt在单点登录系统中的使用
http://blog.leapoahead.com/2015/09/07/user-authentication-with-jwt/
* 使用jwt设计登录系统的流程和思路
http://www.cnblogs.com/binyue/p/4812798.html


# go实现jwt完整例子
------------

## 背景

我的单点登录系统采用go语言开发，在设计中，采用了beego框架。由于beego彼时对session的支持不太好（现在就不得而知了），
又考虑到jwt中可以保存一些用户的基本信息，免去了查库的麻烦，于是就采用了jwt.

但jwt token存在一个用户主动退出登录的问题：
当用户登录时，服务器会返回一个token给客户端（如保存在浏览器中）。因为登录状态（token失效时间）已在token中设置，
该时间无法在当前token串中修改，除非重新生成一个token。所以，当用户主动退出登录时，问题就来了：

1.token是在用户登录时给用户颁发的令牌（用于用户与服务器间认证），用户主动退出时要重新生成一个时间已经失效的token返回给用户么？
是的话，客户端就得用这个新的token替换原来旧的token。但旧token仍未失效，只要该用户访问时带着它，登录认证一样可以通过。

2.如果不重生成token的话，客户端原来保存的token仍在有效期内，此时带着该token访问服务器时，
登录认证是通过的，但前提是用户已经主动注销了登录。这显然是不可以的。

基于这个情况，可以将用户的状态保存在服务器中，这样就不用在token中设置失效时间了。这里我选择用redis来保存，只存储用户的失效时间
，key可以是uid,value是token失效时间，每当用户访问系统时，重置失效时间。
这看起来与session的保存机制是一样的，确实如此。但这不是session，在服务器保存的session是一个比较大的数据结构，
相比一个简单token失效时间占用的内存要大多了。

> 单点登录系统的完整项目代码，请移步我的GitHub：https://github.com/yingyan003/
这里只摘取jwt的部分

## 设计
这里采用jwt的第三方包github.com/dgrijalva/jwt-go来实现。
其中的包含的token的加密算法有很多种，不同的加密算法，需要的参数不一样。
最简单的加密算法是```HS256```

* 采用HS256算法

在生成token时，需要指定一个key，解析时也必须通过同样的key。
```token.SignedString([]byte(key))```
这里的SignedString参数必须是[]byte类型。
> 这是个历史问题，这个第三方包旧版本使用的是string类型，出于安全性考虑，后来改为了[]byte类型。

* 采用ES256算法

生成token时，需要提供一个privateKey:
```token.SignedString(&privateKey)```
生成这个privateKey需要类似下面的3个参数（名称随意）：
```
ECDSAKeyD=CCFDFDC9C2572D15C639D07E3C6C8804A1E941B13F5D10C7297A2DFAA70E6393
ECDSAKeyX=EE4C3E11EB1BF081CFD4B5CCC482E069BFBECA07D566238F29191716319B809E
ECDSAKeyY=A40CCD993EC355326588E2A9E202C24A2D5D1BE5128B19885FD9F2C4155C3EF1
```
这三个参数是用来生成privateKey的，不能乱写，可以用下面的方式生成：
```
    //todo 生成ecdsa.PrivateKey
	randKey:=rand.Reader
	var err error
	prk, err = ecdsa.GenerateKey(elliptic.P256(),randKey )
	if err!=nil{
		fmt.Println("generate key error",err)
	}
	puk=prk.PublicKey
	fmt.Printf("prkD=%X\n",prk.D)
	fmt.Printf("prkX=%X\n",prk.X)
	fmt.Printf("prkY=%X\n",prk.Y)
	fmt.Println("prk",prk," \npbk",puk)
```
其中的prk.D，prk.X,prk.Y对应上面的ECDSAKeyD，ECDSAKeyX,ECDSAKeyY。解析时也会用到

如果采用HS256的方式，由于前端需要解析token获取用户相关的信息，需要把key给到前端。
这样一来，导致key可以随意被获取到。由于前边提到我会把存在redis的用户登录状态的redisKey存在token中，
有了key就可以解析token，这样redisKey就会暴露，这是很不安全的，因为用户登录状态就有被修改的可能。
基于这个原因，我使用了双重token。也就是在一个token里保存另一个token。
最外一层token只保存用户基本且不重要的信息，所以采用HS256算法，前端暴露key也无关紧要。
内嵌的token采用的ES256算法，因为生成和解析token的D/X/Y只保存在服务器中，所以理论上是相对安全的。


## 实现

> 这里用到了jwt的第三方包：github.com/dgrijalva/jwt-go
如果你本地没有的话，可以通过go get拉到本地：
```go get github.com/dgrijalva/jwt-go ```

* 双重token
-> jwt.go

* ES256 demo
-> jwt_test.go


