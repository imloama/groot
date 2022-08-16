package groot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const DEFAULT_CONFIG_FILENAME = "config"
const DEFAULT_CONFIG_FILETYPE = "toml"
const KEY_DEPLOY_ENV = "DEPLOY_ENV"

type ServerConfig struct {
	Host string     `json:"host"`
	Port int        `json:"port"`
	Cors bool       `json:"cors"`
	Log  *LogConfig `json:"log"`
	// OpenApi *OpenApiConfig
	Swagger *SwaggerInfoData `json:"swagger"`
	Redis   *RedisConfig     `json:"redis"`
}
type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}
type SwaggerInfoData struct {
	Enable      bool
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

var serverCfg *ServerConfig

func LoadConfig(files ...string) error {
	if serverCfg != nil {
		return errors.New("init config over")
	}
	serverCfg = &ServerConfig{}
	filename := DEFAULT_CONFIG_FILENAME
	filetype := DEFAULT_CONFIG_FILETYPE
	filepath := "./config"
	if len(files) >= 3 {
		filepath = files[0]
		filename = files[1]
		filetype = files[2]
	} else if len(files) == 2 {
		filepath = files[0]
		filename = files[1]
	} else if len(files) == 1 {
		filename = files[0]
	}

	deployEnv := os.Getenv(KEY_DEPLOY_ENV)
	configFile := ""
	if deployEnv != "" {
		configFile = fmt.Sprintf("%s-%s", filename, deployEnv)
	} else {
		configFile = filename
	}
	fmt.Println("read config file", "configFile", configFile, "filetype", filetype, "filepath", filepath)
	viper.SetConfigName(configFile)
	viper.SetConfigType(filetype)
	viper.AddConfigPath(filepath) // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(err)
	}
	err = viper.UnmarshalKey("server", serverCfg)
	if err != nil {
		fmt.Println("error read config file", err.Error())
		panic(err)
	}
	// fmt.Printf("host = %s\n", viper.GetString("server.host"))
	data, _ := json.Marshal(serverCfg)
	fmt.Println("config file content :", string(data))
	return nil

}
