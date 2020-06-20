package constants

import "time"

var (
	MASTER_IP       = "127.0.0.1"
	MASTER_PORT     = "9527"
	MASTER_MONITER  = 10
	MASTER_NOTIFY   = false
	MASTER_RECEIVER = ""
	MASTER_TICKER   *time.Ticker
	MASTER_TTL      = int64(10)
	MASTER_SCHEMA   = "Asgard"
)
