// config/config.go
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	ES struct {
		URL string
	}
	MongoDB struct {
		URI      string
		Database string
	}
	RabbitMq struct {
		UserName string
		Password string
	}
}

// Load 加载配置信息
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("sql.host", "localhost")
	viper.SetDefault("sql.port", "3306")
	viper.SetDefault("sql.user", "root")
	viper.SetDefault("sql.name", "tdl")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("es.url", "http://localhost:9200")
	viper.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	viper.SetDefault("mongodb.database", "tdl")
	viper.SetDefault("rabbitmq.username", "user")
	viper.SetDefault("rabbitmq.password", "123456")

	// 尝试从环境变量读取配置
	viper.AutomaticEnv()

	// 加载配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在时不返回错误
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetMySQLDSN 返回MySQL数据源名称
func (c *Config) GetMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name)
}
