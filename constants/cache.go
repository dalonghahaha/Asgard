package constants

import "time"

const (
	CACHE_NAME = "asgard"

	CACHE_TTL = 24 * time.Hour

	CACHE_KEY_GROUP = "asgard:group"

	CACHE_KEY_AGENT = "asgard:agent"

	CACHE_KEY_AGENT_IP_PORT = "asgard:agent:ip:port:"

	CACHE_KEY_USER = "asgard:user"
)
