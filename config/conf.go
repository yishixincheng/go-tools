package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func init()  {
	fmt.Println("初始化读取配置文件")
	viper.SetConfigName("application")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read config failed: %v", err)
	}
}
