package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Env        string `mapstructure:"env"`
	LogLevel   string `mapstructure:"log_level"`
	LogFormat  string `mapstructure:"log_format"`
	ServerPort int    `mapstructure:"server_port"`
	DbHost     string `mapstructure:"database_host"`
	DbPort     int    `mapstructure:"database_port"`
	DbUser     string `mapstructure:"database_user"`
	DbPassword string `mapstructure:"database_password"`
	DbName     string `mapstructure:"database_name"`
	JWTSecret  string `mapstructure:"jwt_secret"`
}

func LoadConfig(path string) (*Config, error) {
	// Определяем файл конфигурации
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	// Считывание переменных окружения для переопределения настроек в разных окружениях
	viper.AutomaticEnv()

	// Для корректного именования переменных окружения: database.host -> AVITO_SHOP_DATABASE_HOST
	viper.SetEnvPrefix("AVITO_SHOP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Чтение файла конфигурации
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл конфигурацию: %w", err)
	}

	// config.yaml -> Config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("не удалось декодировать конфигурацию: %w", err)
	}

	return &cfg, nil
}
