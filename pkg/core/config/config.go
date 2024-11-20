package config

import "github.com/spf13/viper"

type Configuration struct {
	MQTTClientName string `mapstructure:"CLIENT_NAME"`
	MQTTURI        string `mapstructure:"URL"`
	MQTTIP         string `mapstructure:"IP"`
	MQTTPort       string `mapstructure:"PORT"`
	MQTTBrokerURL  string `mapstructure:"BROKER_URL"`
}

var (
	GlobalConfig *Configuration
)

func LoadConfig() (*Configuration, error) {

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Configuration

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	setDefaultValues(&cfg)
	GlobalConfig = &cfg
	return &cfg, nil
}

func setDefaultValues(cfg *Configuration) {
	cfg.MQTTBrokerURL = cfg.MQTTURI + cfg.MQTTIP + cfg.MQTTPort
}
