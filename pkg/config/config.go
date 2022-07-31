package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Token      string
	MyID       int
	Delay      int
	DelMin     int
	DelMax     int
	DirParh    string
	AlbumID    int
	MainPostID int
}

func Init() (*Config, error) {
	var cfg Config

	if err := fromEnv(&cfg); err != nil {
		return nil, err
	}

	cfg.Delay = 60000 //17editinPost via 15sec+-0.5 -> captcha  //5editingPost via 30sec+-1sec -> https://api.vk.com/method/wall.edit
	cfg.DelMin = 60000
	cfg.DelMax = 60000 * 2
	cfg.DirParh = "./com/"
	cfg.AlbumID = 235938491
	cfg.MainPostID = 335

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
