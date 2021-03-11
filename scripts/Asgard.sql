
SET NAMES utf8mb4;

CREATE TABLE `agents` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '实例ID',
  `alias` varchar(32) NOT NULL DEFAULT '' COMMENT '别名',
  `ip` varchar(100) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `port` varchar(100) NOT NULL DEFAULT '' COMMENT '端口号',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='实例信息';

CREATE TABLE `apps` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '应用ID',
  `group_id` int(11) NOT NULL DEFAULT 0 COMMENT '分组ID',
  `name` varchar(1000) NOT NULL DEFAULT '' COMMENT '应用ID',
  `agent_id` int(11) NOT NULL DEFAULT 0 COMMENT '运行实例ID',
  `dir` varchar(2000) NOT NULL DEFAULT '' COMMENT '执行目录',
  `program` varchar(2000) NOT NULL DEFAULT '' COMMENT '执行程序',
  `args` varchar(2000) NOT NULL DEFAULT '' COMMENT '执行参数',
  `std_out` varchar(2000) NOT NULL DEFAULT '' COMMENT '标准输出路径',
  `std_err` varchar(2000) NOT NULL DEFAULT '' COMMENT '错误输出路径',
  `auto_restart` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否自动重启',
  `is_monitor` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否开启监控',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `creator` int(11) NOT NULL COMMENT '创建人',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `updator` int(11) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='应用信息';

CREATE TABLE `archives` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `type` tinyint(4) unsigned NOT NULL DEFAULT 0 COMMENT '类别',
  `related_id` tinyint(11) unsigned NOT NULL DEFAULT 0 COMMENT '关联ID',
  `uuid` varchar(100) NOT NULL DEFAULT '' COMMENT '唯一ID',
  `pid` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '运行ID',
  `begin_time` datetime NOT NULL COMMENT '开始运行时间',
  `end_time` datetime NOT NULL COMMENT '结束运行时间',
  `status` tinyint(4) NOT NULL COMMENT '状态',
  `signal` varchar(50) NOT NULL DEFAULT '' COMMENT '信号',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `related` (`type`,`related_id`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='归档信息';

CREATE TABLE `exceptions` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `type` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '类别',
  `related_id` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '关联ID',
  `desc` varchar(1000) NOT NULL DEFAULT '' COMMENT '错误描述',
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `groups` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '分组ID',
  `name` varchar(500) NOT NULL DEFAULT '' COMMENT '分组名称',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `creator` int(11) NOT NULL COMMENT '创建人',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `updator` int(11) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分组信息';

CREATE TABLE `jobs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '任务ID',
  `group_id` int(11) NOT NULL DEFAULT 0 COMMENT '分组ID',
  `name` varchar(1000) NOT NULL DEFAULT '' COMMENT '应用ID',
  `agent_id` int(11) NOT NULL DEFAULT 0 COMMENT '运行实例ID',
  `dir` varchar(2000) NOT NULL DEFAULT '' COMMENT '执行目录',
  `program` varchar(2000) NOT NULL DEFAULT '' COMMENT '执行程序',
  `args` varchar(2000) NOT NULL DEFAULT '' COMMENT '执行参数',
  `std_out` varchar(2000) NOT NULL DEFAULT '' COMMENT '标准输出路径',
  `std_err` varchar(2000) NOT NULL DEFAULT '' COMMENT '错误输出路径',
  `spec` varchar(50) NOT NULL DEFAULT '0' COMMENT '时间配置',
  `timeout` int(11) NOT NULL DEFAULT -1 COMMENT '超时设置(秒)',
  `is_monitor` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否开启监控',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `creator` int(11) NOT NULL COMMENT '创建人',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `updator` int(11) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='计划任务信息';

CREATE TABLE `monitors_202005` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `type` tinyint(4) NOT NULL COMMENT '类别',
  `related_id` int(11) NOT NULL COMMENT '关联ID',
  `uuid` varchar(100) NOT NULL DEFAULT '' COMMENT '唯一ID',
  `pid` int(11) NOT NULL COMMENT '运行ID',
  `cpu` float NOT NULL COMMENT 'CPU占用',
  `memory` float NOT NULL COMMENT '内存占用',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `monitors_202006` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `type` tinyint(4) NOT NULL COMMENT '类别',
  `related_id` int(11) NOT NULL COMMENT '关联ID',
  `uuid` varchar(100) NOT NULL DEFAULT '' COMMENT '唯一ID',
  `pid` int(11) NOT NULL COMMENT '运行ID',
  `cpu` float NOT NULL COMMENT 'CPU占用',
  `memory` float NOT NULL COMMENT '内存占用',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `monitors_202007` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `type` tinyint(4) NOT NULL COMMENT '类别',
  `related_id` int(11) NOT NULL COMMENT '关联ID',
  `uuid` varchar(100) NOT NULL DEFAULT '' COMMENT '唯一ID',
  `pid` int(11) NOT NULL COMMENT '运行ID',
  `cpu` float NOT NULL COMMENT 'CPU占用',
  `memory` float NOT NULL COMMENT '内存占用',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `operations` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '操作人',
  `type` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '关联对象类别',
  `related_id` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '关联对象ID',
  `action` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '动作',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `timings` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '任务ID',
  `group_id` int(11) NOT NULL DEFAULT 0 COMMENT '分组ID',
  `name` varchar(1000) NOT NULL DEFAULT '' COMMENT '应用ID',
  `agent_id` int(11) NOT NULL DEFAULT 0 COMMENT '运行实例ID',
  `dir` varchar(2000) NOT NULL DEFAULT '' COMMENT '执行目录',
  `program` varchar(2000) NOT NULL DEFAULT '' COMMENT '执行程序',
  `args` varchar(2000) NOT NULL DEFAULT '' COMMENT '执行参数',
  `std_out` varchar(2000) NOT NULL DEFAULT '' COMMENT '标准输出路径',
  `std_err` varchar(2000) NOT NULL DEFAULT '' COMMENT '错误输出路径',
  `time` datetime NOT NULL COMMENT '执行时间',
  `timeout` int(11) NOT NULL DEFAULT -1 COMMENT '超时设置(秒)',
  `is_monitor` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否开启监控',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `creator` int(11) NOT NULL COMMENT '创建人',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `updator` int(11) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='定时任务信息';

CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `nickname` varchar(100) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(500) NOT NULL DEFAULT '' COMMENT '头像',
  `email` varchar(200) NOT NULL DEFAULT '' COMMENT '邮箱',
  `mobile` varchar(50) NOT NULL DEFAULT '' COMMENT '手机号',
  `role` varchar(32) NOT NULL DEFAULT '' COMMENT '角色',
  `salt` varchar(50) NOT NULL DEFAULT '' COMMENT '加密秘钥',
  `password` varchar(200) NOT NULL DEFAULT '' COMMENT '密码',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息';

INSERT INTO `users` (`id`, `nickname`, `avatar`, `email`, `mobile`, `role`, `salt`, `password`, `status`, `created_at`, `updated_at`)
VALUES	(1, 'admin', '', 'admin@admin.com', '13888888888', 'Administrator', 'PWxWFtzI', 'd55d350c7d70b3f5fa295ca05f98af73', 1, '2020-01-08 16:38:05', '2020-05-20 08:01:57');

