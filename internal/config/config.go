package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/k1v4/avito_shop/pkg/DB/postgres"
	"github.com/k1v4/avito_shop/pkg/DB/redis"
)

type Config struct {
	postgres.DBConfig
	redis.RedisConfig

	RestServerPort int `env:"REST_SERVER_PORT" env-description:"rest server port" env-default:"8080"`
}

func MustLoadConfig() *Config {
	//errEnv := godotenv.Load(".env")
	//if errEnv != nil {
	//	panic(errEnv)
	//}

	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
