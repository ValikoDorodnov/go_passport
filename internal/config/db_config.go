package config

type DbConfig struct {
	User    string `mapstructure:"db_user"`
	Pass    string `mapstructure:"db_password"`
	Name    string `mapstructure:"db_name"`
	Port    string `mapstructure:"db_port"`
	Host    string `mapstructure:"db_host"`
	SllMode string `mapstructure:"db_ssl"`
}
