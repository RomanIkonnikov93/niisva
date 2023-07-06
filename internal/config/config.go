package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	GRPCAddress     string `env:"SERVER_ADDRESS" envDefault:":3200"`
	FileStoragePath string `env:"STORAGE_PATH" envDefault:"storage"`
	FileType        string `env:"FILE_TYPE" envDefault:".txt"`
	UsersPathStore  string `env:"USERS_PATH_STORE" envDefault:"users.csv"` // .csv file
}

func GetConfig() (*Config, error) {

	cfg := &Config{}

	flag.StringVar(&cfg.GRPCAddress, "g", cfg.GRPCAddress, "SERVER_ADDRESS")
	flag.StringVar(&cfg.FileStoragePath, "s", cfg.FileStoragePath, "FILE_STORAGE_PATH")
	flag.StringVar(&cfg.FileType, "f", cfg.FileType, "FILE_TYPE")
	flag.StringVar(&cfg.UsersPathStore, "u", cfg.UsersPathStore, "USERS_PATH_STORE")
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
