package config

type JwtConfig struct {
	AccessTtl  int    `mapstructure:"jwt_access_ttl"`
	RefreshTtl int    `mapstructure:"jwt_refresh_ttl"`
	SecretKey  string `mapstructure:"jwt_secret_key"`
	Issuer     string `mapstructure:"jwt_issuer"`
}
