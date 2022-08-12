package groot

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const DEFAULT_CONFIG_FILENAME = "config"
const DEFAULT_CONFIG_FILETYPE = "toml"
const KEY_DEPLOY_ENV = "DEPLOY_ENV"

func LoadConfig(files ...string) error {
	filename := DEFAULT_CONFIG_FILENAME
	filetype := DEFAULT_CONFIG_FILETYPE
	filepath := "."
	if len(files) >= 3 {
		filepath = files[0]
		filename = files[1]
		filetype = files[2]
	} else if len(files) == 2 {
		filepath = files[0]
		filename = files[1]
	} else {
		filename = files[0]
	}

	deployEnv := os.Getenv(KEY_DEPLOY_ENV)
	configFile := ""
	if deployEnv != "" {
		configFile = fmt.Sprintf("%s-%s", filename, deployEnv)
	}
	viper.SetConfigName(configFile)
	viper.SetConfigType(filetype)
	viper.AddConfigPath(filepath) // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return err
	}
	return nil

}

func UnmarshalConfig(data interface{}) error {
	return viper.Unmarshal(data)
}

func UnmarshalConfigByKey(key string, data interface{}) error {
	return viper.UnmarshalKey(key, data)
}
