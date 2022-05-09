package config

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
	db "tradingdata/internal/db/sqlc"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	ClientId    string `mapstructure:"CLIENT_ID"`
	SecretKey   string `mapstructure:"SECRET_KEY"`
	Pan         string `mapstructure:"PAN"`
	UserName    string `mapstructure:"USERNAME"`
	Password    string `mapstructure:"PASSWORD"`
	RedirectUrl string `mapstructure:"REDIRECT_URL"`
	DbHostname  string `mapstructure:"DB_HOSTNAME"`
	DbPort      int64  `mapstructure:"DB_PORT"`
	DbUsername  string `mapstructure:"DB_USERNAME"`
	DbPassword  string `mapstructure:"DB_PASSWORD"`
	DbDatabase  string `mapstructure:"DB_DATABASE"`
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
}

// global instance ,should be reused.
type GlobalInstance struct {
	Config  Config
	TokenDb db.TokenDB
}

var masterConfig Config

var masterConfigOnce sync.Once

// initialize env config
func GetConfig(path string, filename string, fileextension string) Config {
	masterConfigOnce.Do(func() {
		masterConfig = LoadConfig(path, filename, fileextension, Config{}).(Config)
	})
	return masterConfig
}

// seting up the job db.
func SetupDb(config Config) (*sql.DB, error) {
	psqlInfo := config.DatabaseUrl
	if config.Environment == "development" {
		psqlInfo = fmt.Sprintf(PostgresSource,
			config.DbHostname, config.DbPort, config.DbUsername, config.DbPassword, config.DbDatabase)
	}
	return sql.Open(DbDriver, psqlInfo)
}

func LoadConfig(path string, filename string, fileextension string, confg interface{}) (config interface{}) {
	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType(fileextension)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&confg)
	if err != nil {
		log.Fatal(err)
	}
	return confg
}
