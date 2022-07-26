package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Token   string
	MyID    int
	Delay   int
	DirParh string
}

func Init() (*Config, error) {
	var cfg Config

	if err := fromEnv(&cfg); err != nil {
		return nil, err
	}

	cfg.Delay = 60000 * 3
	cfg.DirParh = "./assets/"

	return &cfg, nil
}

func fromEnv(cfg *Config) error {
	godotenv.Load()

	if err := viper.BindEnv("token"); err != nil {
		return err
	}
	cfg.Token = viper.GetString("token")

	if err := viper.BindEnv("my_id"); err != nil {
		return err
	}
	cfg.MyID = viper.GetInt("my_id")

	return nil
}
