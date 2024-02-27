package config

import (
	"github.com/BelyaevEI/wallet/internal/models"
	"github.com/spf13/viper"
)

type Config struct {
	DSN  string `mapstructure:"DSN"`  // DSN for postgreSQL
	Host string `mapstructure:"Host"` // Server host
	Port string `mapstructure:"Port"` // Server port
}

// Reading config file for setting application
func LoadConfig(path string) (Config, error) {

	conf := Config{}

	viper.AddConfigPath(path)
	viper.SetConfigName(models.ConfigName)
	viper.SetConfigType(models.ConfigType)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return conf, err
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}
