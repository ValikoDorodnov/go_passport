package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	envKeysMap := map[string]interface{}{}

	if err := mapstructure.Decode(rawVal, &envKeysMap); err != nil {
		return err
	}

	for k := range envKeysMap {
		if bindErr := viper.BindEnv(k); bindErr != nil {
			return bindErr
		}
	}

	return viper.Unmarshal(rawVal, opts...)
}
