package app

import (
	"log"

	"github.com/YeHeng/gtool/common/model"

	"github.com/spf13/viper"
)

var Config model.Configuration

func LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./etc/")
	viper.AddConfigPath("/etc/gtool")
	viper.AddConfigPath("$HOME/.gtool")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

}
