package cm

const (
	AgentConfigKey = "config.json"
)

func GenZkAgentConfig() (string, error) {
	var config string
	config = `{
    "debug": true,
    "zkHost": "127.0.0.1",
    "zkPort": "2181",
    "http": {
        "enabled": true,
        "listen": ":1988",
        "backdoor": true
    }
}`
	return config, nil
}
