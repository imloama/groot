package groot

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const DEFAULT_CONFIG_FILENAME = "config"
const DEFAULT_CONFIG_FILETYPE = "toml"
const KEY_DEPLOY_ENV = "DEPLOY_ENV"

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Cors bool   `json:"cors"`
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
	JsonFile    string // json file path
}

var serverCfg = ServerConfig{}
var logCfg = LogConfig{}
var swaggerCfg = SwaggerInfoData{}
var redisCfg = RedisConfig{}
var ormCfg = OrmConfig{}

func LoadConfig(files ...string) error {
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
	err = viper.UnmarshalKey("server", &serverCfg)
	if err != nil {
		fmt.Println("error read config file", err.Error())
		panic(err)
	}
	viper.UnmarshalKey("log", &logCfg)
	viper.UnmarshalKey("swagger", &swaggerCfg)
	viper.UnmarshalKey("redis", &redisCfg)
	viper.UnmarshalKey("orm", &ormCfg)
	return nil

}

func GetSwaggerConfig() SwaggerInfoData {
	return swaggerCfg
}

func GetServerConfig() ServerConfig {
	return serverCfg
}
func GetLogConfig() LogConfig {
	return logCfg
}
func GetRedisConfig() RedisConfig {
	return redisCfg
}
func GetOrmConfig() OrmConfig {
	return ormCfg
}
