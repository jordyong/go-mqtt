package config

type mqttConfig struct {
	ClientName string
	URI        string
	IP         string
	Path       string
	BrokerURL  string
}

type Configuration struct {
	MQTTopts mqttConfig
}

var (
	GlobalConfig *Configuration
)

func LoadConfig() (*Configuration, error) {
	var cfg Configuration
	setDefaultValues(&cfg)
	GlobalConfig = &cfg
	return &cfg, nil
}

func setDefaultValues(cfg *Configuration) {
	cfg.MQTTopts.ClientName = "go-mqtt"
	cfg.MQTTopts.URI = "ws://"
	cfg.MQTTopts.IP = "172.27.2.35"
	cfg.MQTTopts.Path = ":9001/mqtt"
	cfg.MQTTopts.BrokerURL = cfg.MQTTopts.URI + cfg.MQTTopts.IP + cfg.MQTTopts.Path
}
