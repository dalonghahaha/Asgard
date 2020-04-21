package constants

import "github.com/dalonghahaha/avenger/structs"

const (
	TYPE_AGENT  = int64(1)
	TYPE_APP    = int64(2)
	TYPE_JOB    = int64(3)
	TYPE_TIMING = int64(4)

	AGENT_OFFLINE   = int64(0)
	AGENT_ONLINE    = int64(1)
	AGENT_FORBIDDEN = int64(-1)

	GROUP_STATUS_UNUSAGE = int64(0)
	GROUP_STATUS_USAGE   = int64(1)

	APP_STATUS_DELETED = int64(-1)
	APP_STATUS_STOP    = int64(0)
	APP_STATUS_RUNNING = int64(1)
	APP_STATUS_PAUSE   = int64(2)

	JOB_STATUS_DELETED = int64(-1)
	JOB_STATUS_STOP    = int64(0)
	JOB_STATUS_RUNNING = int64(1)
	JOB_STATUS_PAUSE   = int64(2)

	TIMING_STATUS_DELETED  = int64(-1)
	TIMING_STATUS_STOP     = int64(0)
	TIMING_STATUS_RUNNING  = int64(1)
	TIMING_STATUS_PAUSE    = int64(2)
	TIMING_STATUS_FINISHED = int64(3)
)

var APP_STATUS = []structs.M{
	{
		"ID":   APP_STATUS_STOP,
		"Name": "停止",
	},
	{
		"ID":   APP_STATUS_RUNNING,
		"Name": "运行中",
	},
	{
		"ID":   APP_STATUS_PAUSE,
		"Name": "暂停",
	},
	{
		"ID":   APP_STATUS_DELETED,
		"Name": "已删除",
	},
}

var JOB_STATUS = []structs.M{
	{
		"ID":   JOB_STATUS_STOP,
		"Name": "停止",
	},
	{
		"ID":   JOB_STATUS_RUNNING,
		"Name": "运行中",
	},
	{
		"ID":   JOB_STATUS_PAUSE,
		"Name": "暂停",
	},
	{
		"ID":   JOB_STATUS_DELETED,
		"Name": "已删除",
	},
}

var TIMING_STATUS = []structs.M{
	{
		"ID":   TIMING_STATUS_STOP,
		"Name": "停止",
	},
	{
		"ID":   TIMING_STATUS_RUNNING,
		"Name": "运行中",
	},
	{
		"ID":   TIMING_STATUS_PAUSE,
		"Name": "暂停",
	},
	{
		"ID":   TIMING_STATUS_FINISHED,
		"Name": "已完成",
	},
	{
		"ID":   TIMING_STATUS_DELETED,
		"Name": "已删除",
	},
}
