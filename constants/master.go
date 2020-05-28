package constants

import "time"

var (
	MASTER_IP             = "127.0.0.1"
	MASTER_PORT           = 9527
	MASTER_TIMEOUT        = time.Second * 10
	MASTER_REPORT_TIMEOUT = time.Second * 5
	MASTER_MONITER        = 10
	MASTER_TICKER         *time.Ticker
)
