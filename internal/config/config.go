package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config содержит конфигурацию приложения.
type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	StoragePath    string     `yaml:"storage_path" env-required:"true"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

// GRPCConfig содержит конфигурацию gRPC сервера.
type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// MustLoad загружает конфигурацию из файла и переменных окружения.
func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath возвращает путь к файлу конфигурации.
func fetchConfigPath() string {
	var res string
	// Проверяем флаги командной строки
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	// Если флаг не установлен, смотрим переменную окружения
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
