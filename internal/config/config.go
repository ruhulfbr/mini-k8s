package config

type Config struct {
	App    AppConfig
	HTTP   HTTPConfig
	Redis  RedisConfig
	Badger BadgerConfig
	Docker DockerConfig
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
	Port int `env:"REDIS_PORT" envDefault:"6379"`
}

type BadgerConfig struct {
	DataSource string `env:"BADGER_DATA_SOURCE" envDefault:"data/badger"`
}

type DockerConfig struct {
	ContainerNamePref string `env:"CONTAINER_NAME_PREFIX"`
	ImageTagPrefix    string `env:"IMAGE_TAG_PREFIX"`
}
