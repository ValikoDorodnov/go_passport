package config

type RedisConfig struct {
	Pass string `mapstructure:"redis_password"`
	Port string `mapstructure:"redis_port"`
	Host string `mapstructure:"redis_host"`
	Db   int    `mapstructure:"redis_db"`
}
