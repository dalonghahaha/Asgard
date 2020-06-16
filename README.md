# Asgard

## 简介

Asgard是设计用于综合解决常驻进程应用、计划任务、定时任务的分布式作业管理系统。

## 架构设计

![架构设计图](/doc/Asgard.png)

- Asgard系统由web节点、master节点、agent节点组成。
- web节点主要功能包括实例管理、分组管理、作业配置、作业运行状态控制、作业运行状态查看、日志查询
- master节点负责agent节点的状态监测，同时接收并转存agent节点上报的运行时数据
- agent节点接收web节点的指令在相应的服务器中运作相关作业
- master节点和agent节点之间通过grpc协议交换数据

## Web界面预览

![首页控制台](/doc/page1.png)

![实例管理](/doc/page2.png)

![任务管理](/doc/page3.png)

![监控信息](/doc/page4.png)

## 指令作用

### Asgard guard

启动为管理单机系统常驻进程应用守护服务，类似supervisor

#### Asgard guard status

查看单机常驻进程应用列表

#### Asgard guard show xxxxxx

查看单机常驻进程应用xxxxxx信息

### Asgard cron

启动为管理单机系统的计划任务守护服务，类似crontab

#### Asgard cron status

查看单机计划任务列表

#### Asgard cron show xxxxxx

查看单机计划任务xxxxxx信息

### Asgard web

启动为web节点

### Asgard msater

启动为master节点

### Asgard agent

启动为agent节点

### Asgard agent status

查看agent节点运行的常驻进程应用、计划任务、定时任务综合信息

## 部署及配置说明

运行**scripts/Asgard.sql**初始化mysql数据表。monitors数据表按月分表，格式"monitors_202006",需要提前创建。

web节点、master节点、agent节点都需要一个名为**app.yaml**的配置文件。

默认读取运行目录下名为**conf**的目录，如果需要指向配置文件的目录可以在启动命令中通过**conf**参数指定。

启动web节点需要用到web目录中的**assets**和**views**两部分静态资源

### master节点配置项示例及说明

``` yaml
master:
    port: 9527                          #master节点监听端口
component:
    db:
        asgard:                         #mysql数据库配置
            host: "127.0.0.1"           #mysql数据库地址
            port: 3306                  #mysql数据库端口
            user: "xxxxxx"              #mysql数据库用户名
            password: "xxxxxx"          #mysql数据库密码
            database: "Asgard"          #mysql数据库库名
    redis:                              #redis配置
        asgard:
            host: "127.0.0.1"           #redis地址
            port: 6379                  #redis端口号
            password: ""                #redis密码
            db: 0                       #redis库索引
    log:                                #日志配置
        console: true                   #是否输出到控制台
        level: "debug"                  #日志级别
        dir: "runtime/"                 #日志存放根目录
```

### web节点配置项示例及说明

``` yaml
web:
    port: 12345                         #web节点监听端口
    domain: "asgard.dalong.me"          #web节点域名
    cookie_salt: "sdswqeqx"             #web节点身份验证加密值(必须为8位字符串)
component:
    db:
        asgard:                         #mysql数据库配置
            host: "127.0.0.1"           #mysql数据库地址
            port: 3306                  #mysql数据库端口
            user: "xxxxxx"              #mysql数据库用户名
            password: "xxxxxx"          #mysql数据库密码
            database: "Asgard"          #mysql数据库库名
    redis:                              #redis配置
        asgard:
            host: "127.0.0.1"           #redis地址
            port: 6379                  #redis端口号
            password: ""                #redis密码
            db: 0                       #redis库索引
    log:                                #日志配置
        console: true                   #是否输出到控制台
        level: "debug"                  #日志级别
        dir: "runtime/"                 #日志存放根目录
```

### agent节点配置项示例及说明

``` yaml
system:
    moniter: 10                         #监控指标上报周期，单位秒
    timer: 10                           #定时任务轮询周期，单位秒
agent:
    master:                             #master节点地址
        ip: "127.0.0.1"                 #master节点端口
        port: 9527
    rpc:
        ip: "127.0.0.1"                 #agent节点地址
        port: 27149                     #agent节点端口
component:
    log:                                #日志配置
        console: true                   #是否输出到控制台
        level: "debug"                  #日志级别
        dir: "runtime/"                 #日志存放根目录
```
