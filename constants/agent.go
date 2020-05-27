package constants

import "time"

var (
	AGENT_IP             = ""
	AGENT_PORT           = "27149"
	AGENT_PID            = 0
	AGENT_UUID           = ""
	AGENT_MONITER        = 30
	AGENT_MONITER_TICKER *time.Ticker
)
