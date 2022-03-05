package configmanager

import (
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

//GlobalConfig 全ser可使用的參數
var GlobalConfig *Configuration

// LogSetting 是 log 的設定
type LogSetting struct {
	Name             string `yaml:"name"`
	Type             string `yaml:"type"`
	MinLevel         string `yaml:"min_level"`
	ConnectionString string `yaml:"connection_string"`
}

// Database 用來提供連線的資料庫數據
type Database struct {
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	DBName   string `yaml:"dbname"`
}

type Mongo struct {
	Address  string   `yaml:"address"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Database []string `yaml:"database"`
}
type Mysql struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func init() {
	GlobalConfig = Reload()
}

//Configuration 參數結構
type Configuration struct {
	//WSServer jiface.IServer
	Env string `yaml:"env"`

	Databases []Database `yaml:"databases"`

	SyncDataPlatform struct {
		Platform       string `yaml:"platform"`
		HTTPBind       string `yaml:"http_bind"`
		NotifierCenter struct {
			RabbitMQAddr string `yaml:"rabbit_mq_addr"`
			VirtualHosts string `yaml:"virtual_hosts"`
			Username     string `yaml:"username"`
			Password     string `yaml:"password"`
		} `yaml:"notifier_center"`
		Mysql Mysql `yaml:"mysql"`
	} `yaml:"syncData_platform"`
}

//Reload 重新載入參數
func Reload() *Configuration {
	data, err := ioutil.ReadFile("configs/config.yml")
	if err != nil {
		log.Panic().Msgf("%v", err)
	}

	tempPara := &Configuration{}

	err = yaml.Unmarshal(data, &tempPara)
	if err != nil {
		log.Panic().Msgf("%v", err)
	}

	if value, ok := os.LookupEnv("PLATFORM"); ok {
		tempPara.SyncDataPlatform.Platform = value
	}
	if value, ok := os.LookupEnv("RABBITMQ_ADDRESS"); ok {
		tempPara.SyncDataPlatform.NotifierCenter.RabbitMQAddr = value
	}
	if value, ok := os.LookupEnv("VIRTUAL_HOSTS"); ok {
		tempPara.SyncDataPlatform.NotifierCenter.VirtualHosts = value
	}
	if value, ok := os.LookupEnv("RABBITMQ_USER"); ok {
		tempPara.SyncDataPlatform.NotifierCenter.Username = value
	}
	if value, ok := os.LookupEnv("RABBITMQ_PASSWORD"); ok {
		tempPara.SyncDataPlatform.NotifierCenter.Password = value
	}
	if value, ok := os.LookupEnv("MYSQL_ADDRESS"); ok {
		tempPara.SyncDataPlatform.Mysql.Address = value
	}
	if value, ok := os.LookupEnv("MYSQL_USER"); ok {
		tempPara.SyncDataPlatform.Mysql.Username = value
	}
	if value, ok := os.LookupEnv("MYSQL_PASSWORD"); ok {
		tempPara.SyncDataPlatform.Mysql.Password = value
	}
	if value, ok := os.LookupEnv("MYSQL_DATABASE"); ok {
		tempPara.SyncDataPlatform.Platform = value
	}
	log.Info().Msgf("%v", tempPara)

	return tempPara
}
