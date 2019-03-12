urlooker-web SSO-CAS 版本
============

Web组件主要用于

- 监控项的添加
- 告警组人员管理
- 查看url访问质量/TCP端口连通性绘图

## 新增功能
- 新增SSO Case登陆方式
- 修改收件人添加逻辑，不与用户名强关联，只需要填写邮箱即可，可以任意填写
- SSO登入登出
- 发送到falcon采用直接推送transfer，而不是之前的通过本机的agent http口二次转发给transfer，并且增加冗余
- 增加tcp 端口扫描
- 使用logrus管理日志

## Installation

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/urlooker
cd $GOPATH/src/github.com/urlooker
git clone https://github.com/peng19940915/web.git
cd web
./control build
./control start
```

## Configuration

```
    {
        "log":{
            "logLevel": "debug", # 日志级别
            "logPath": "logs", # 日志文件存放路径
            "fileName": "alarm.log",# 日志文件名字
            "maxAge": 7, # 日志文件存在的最长时间
            "rotationTime": 1 # 日志文件切割周期
        },
        "sso": {
            "serverUrl": "https://sso-cas.xxxdissector.com",# sso server地址
            "serviceUrl": "http://10.139.14.143:1984/auth/login"# 要被代理的地址只需改掉该地址即可
        },
        "admins":["leiyupeng","admin"], # 哪些账号是管理员账号
        "salt": "have fun!",
        "past": 30, # 历史数据展示时间范围
        "http": {
            "listen": "0.0.0.0:1984",# 开放的http端口
            "secret": "secret"
        },
        "rpc": {
            "listen": "0.0.0.0:1985"# RPC端口号
        },
        "mysql": {
            "addr": "root:test456@tcp(10.201.82.28:3306)/urlooker?charset=utf8&&loc=Asia%2FShanghai",
            "idle": 10,
            "max": 20
        },
        "alarm":{
            "enable": true,
            "batch": 200,
            "replicas": 500,
            "connTimeout": 1000,
            "callTimeout": 5000,
            "maxConns": 32,
            "maxIdle": 32,
            "sleepTime":30,
            "cluster":{
                "node-1":"127.0.0.1:1986"
            }
        },
        "monitorMap": {
            "default":["hostname.1"], #监控指标多了之后agent地址可以填多个
            "idc1":["hostname.2"]
        },
        "falcon":{ # Falcon地址
            "enable": true,
            "addrs": [
                "10.203.40.143:8433", # transfer地址，如果有多个，则放置多个
                "10.200.200.71:8433",
                "10.201.81.120:8433",
                "10.202.80.139:8433"
            ],
            "timeout": 3000,
            "interval": 60
        },
        "internalDns":{ #通过公司内部接口获取url对应ip所在机房
            "enable": false,
            "addr":""
        }
    }
    

```

## 问题与使用细节
### 发送请求间隔时间
可以在agent组件中配置：web.interval
### 启动顺序
web > agent > alarm
### agent 模式
由于原版本只支持url 扫描，但这个版本的agent可以同时支持tcp 端口与url拨测  
agent有两种模式  
- port
如果为port 则该agent为tcp port拨测
- url
如果是url 则该agent为url拨测agent

### 如何配置多机房
多机房的意义就在于 能根据域名/主机名解析出来的信息，将条目下发到指定机房的指定agent  
但是问题在于各个公司的cmdb都不太一样，所以该接口预留出来了  
要实现的逻辑为：  
传入域名：返回该域名的IP，以及每个IP所在的机房  
如果没有实现该接口，那所有host默认的机房为default   
default":["hostname.1"] 该段配置的意思hostname.1 这个agent属于default这个机房  
同时机房内如果有多个agent，他的分配机制为：  
将策略的`sid%该机房的agent个数` 会得出一个下标。用该下标去agent 的list内提取出一个agent，该条策略就会被分配到该agent上监控
### agent配置
agent 的主机名，有两种方式生成，第一种为取hostname，第二种手工指定，不管是哪种，都必须存在web组件的配置文件中（monitorMap）
### 日志
日志的存取使用了logrus这一工具，目前默认是存在本地的，并且会做切割，也可以接入到es等工具中
