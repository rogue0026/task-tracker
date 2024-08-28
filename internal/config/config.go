package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type BotCfg struct {
	Token string `yaml:"token"`
}

func Load(path string) (BotCfg, error) {
	_, err := os.Stat(path)
	if err != nil {
		return BotCfg{}, err
	}
	cfg := BotCfg{}
	if err = cleanenv.ReadConfig(path, &cfg); err != nil {
		return BotCfg{}, err
	}
	return cfg, nil
}
