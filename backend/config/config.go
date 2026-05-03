package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitConfig() error {
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(execDir + "/config")
	viper.AddConfigPath("./config")

	return viper.ReadInConfig()
}
