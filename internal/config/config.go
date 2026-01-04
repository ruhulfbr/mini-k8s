package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var (
	cfg  *Config
	once sync.Once
)

type Config struct {
	App    AppConfig
	HTTP   HTTPConfig
	Redis  RedisConfig
	SQLite SQLiteConfig
	Docker DockerConfig
	Logger LogConfig
}

type AppConfig struct {
	Name      string `env:"APP_NAME"`
	PortRange string `env:"APP_PORT_RANGE"`
}

type HTTPConfig struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"API_PORT" envDefault:"3000"`
}

type RedisConfig struct {
	Host string `env:"REDIS_HOST" envDefault:"localhost"`
	Port int    `env:"REDIS_PORT" envDefault:"6379"`
}

type SQLiteConfig struct {
	DataSource string `env:"SQLITE_DATA_SOURCE" envDefault:"data/sqlite/mini-k8s.db"`
}

type DockerConfig struct {
	ContainerNamePref string `env:"CONTAINER_NAME_PREFIX"`
	ImageTagPrefix    string `env:"IMAGE_TAG_PREFIX"`
}

type LogConfig struct {
	Application string `env:"LOG_APPLICATION"`

	// File represents path to file where store logs. Used [os.Stdout] if empty.
	File string `env:"LOG_FILE"`

	// One of: "DEBUG", "INFO", "WARN", "ERROR". Default: "DEBUG".
	Level string `env:"LOG_LEVEL" envDefault:"DEBUG"`

	// Add source code position to messages.
	AddSource bool `env:"LOG_ADD_SOURCE"`
}

func Load() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		var c Config
		if err := env.Parse(&c); err != nil {
			log.Fatalf("load config: %v", err)
		}
		cfg = &c
	})

	return cfg
}
