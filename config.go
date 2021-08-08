package main

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ConfigParser struct {
	Telegram Telegram
	MysqlDB  Database
}

type Telegram struct {
	BotKey string
}

type Database struct {
	Host            string
	Port            int
	Username        string
	Password        string
	DBName          string
	Charset         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int64
}

// 加载配置
func WatchConfig(changeConfig chan struct{}, filename string) error {

	AppPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	viper.AddConfigPath(filepath.Join(AppPath, "config"))
	viper.SetConfigName(filename)
	viper.SetConfigType("toml")

	if err = viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		changeConfig <- struct{}{}
	})

	return nil
}

// 重新从viper载入与校验
func LoadConfig() (*ConfigParser, error) {

	var config ConfigParser
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
