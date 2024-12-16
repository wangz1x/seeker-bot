// load - 2024/12/16
// Author: wangzx
// Description:

package conf

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log/slog"
)

var GvaConfig *Setting

func Load() {
	// name of config file (without extension)
	viper.SetConfigName("setting")
	// REQUIRED if the config file does not have the extension in the name
	viper.SetConfigType("yaml")
	// optionally look for config in the working directory
	viper.AddConfigPath(".")
	// optionally look for config in the working directory
	viper.AddConfigPath("./conf")
	// Find and read the config file
	err := viper.ReadInConfig()
	// Handle errors reading the config file
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	// unmarshal settings
	if err := viper.Unmarshal(&GvaConfig); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	// validate
	validate := validator.New()
	if err := validate.Struct(GvaConfig); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	slog.Info("config loaded", "config_path", viper.ConfigFileUsed())
}

func init() {
	Load()
}
