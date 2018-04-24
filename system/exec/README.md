# go实现shell命令

> 背景

接触系统底层开发的程序员，有时需要用程序实现系统命令，如ping,telnet等等。

这里给大家介绍用go如果实现系统命令的执行。场景就像在本机执行shell命令一样。

### 实现说明

* 关键包：exec
go本身为我们提供了一个exec包。有了它，用go执行系统命令变得so easy。只需1，2行代码就能搞定。

* 关键代码
```
cmd := exec.Command("/bin/bash", "-c", "netstat -ant"+grepPort)
err := cmd.Run()
```

* 执行单条命令
```
cmd := exec.Command("ping", "-c", count, dstIP)
```
> 使用说明：
第一个参数：命令名
剩余参数：命令参数
每个参数都一般来说都要用"，"隔开

错误例子:
```
cmd := exec.Command("netstat", "-ant","|","grep","xxx")
```
因为grep本身也是一个命令，这样用不会报错，但是grep并不能筛选出指定信息。
一般来说会返回netstat -ant所有结果，但有时又只返回部分不知什么规则的信息，比较奇怪。

**所以，执行单条命令时，最好不要掺杂其他命令。否则会得不到期待的结果**

* 执行多条命令

```
cmd := exec.Command("/bin/bash", "-c", "命令")
//如：
cmd := exec.Command("/bin/bash", "-c", "netstat -ant | grep 443")
```
> 使用说明
" /bin/bash -c 命令 "是一般系统（linux,window,mac等）都支持的shell执行命令格式。
>
/bin/bash:
使用的是/bin/bash shell来执行命令。
当然，替换成/bin/sh或bash或sh也行。取决于你的系统执行命令的shell有哪些。
为何直接输入sh或bash也行呢？那是因为linux系统一般都已将/bin添加到系统环境变量中，所以当只输入sh时，系统可以找到/bin/sh来执行命令。
>
-c:
-c是固定的参数


### 实现
这里以netstat为例子，给出了实现demo
