package cm

const (
	AgentConfigKey = "config.json"
)

func GenZkAgentConfig() string {
	return `{
    "debug": true,
    "zkHost": "127.0.0.1",
    "zkPort": "2181",
    "http": {
        "enabled": true,
        "listen": ":1988",
        "backdoor": true
    }
}`
}
