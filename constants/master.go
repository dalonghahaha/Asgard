package constants

import "time"

var (
	MASTER_IP                   = "127.0.0.1"
	MASTER_PORT                 = "9527"
	MASTER_MONITER              = 10
	MASTER_NOTIFY               = false
	MASTER_RECEIVER             = ""
	MASTER_TICKER               *time.Ticker
	MASTER_CLUSTER              = false
	MASTER_CLUSTER_REGISTRY     = []string{}
	MASTER_CLUSTER_SCHEMA       = "Asgard"
	MASTER_CLUSTER_NAME         = ""
	MASTER_CLUSTER_ID           = ""
	MASTER_CLUSTER_IP           = ""
	MASTER_CLUSTER_TTL          = int64(10)
	MASTER_CLUSTER_CAMPAIGN_TTL = 30
)
