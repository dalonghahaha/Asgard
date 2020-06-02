# Asgard

## 简介

Asgard是设计用于综合解决常驻进程应用、计划任务、定时任务的分布式作业管理系统。

## 架构设计

- Asgard系统由web节点、master节点、agent节点组成。
- web节点主要功能包括实例管理、分组管理、作业配置、作业运行状态控制、作业运行状态查看
- master节点负责agent节点的状态监测，同时接收并转存agent节点上报的运行时数据
- agent节点接收web节点的指令在相应的服务器中运作相关作业
- master节点和agent节点之间通过grpc协议交换数据

## 指令作用

### Asgard guard

启动为管理单机系统常驻进程应用守护服务，类似supervisor

### Asgard cron

启动为管理单机系统的计划任务守护服务，类似crontab

### Asgard web

启动为web节点

### Asgard msater

启动为master节点

### Asgard agent

启动为agent节点

### Asgard agent status

查看agent节点运行状态
