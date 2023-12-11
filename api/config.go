package api

import "github.com/spf13/viper"

type Config interface {
	GetAppPort() int
	GetRedisAddress() string
	GetRedisPass() string
	GetRedisDBName() int
}

type config struct {
	env *viper.Viper
}

func (c *config) GetAppPort() int {
	c.env.AutomaticEnv()
	return c.env.GetInt("app_port")
}

func (c *config) GetRedisAddress() string {
	c.env.AutomaticEnv()
	return c.env.GetString("redis_address")
}

func (c *config) GetRedisPass() string {
	c.env.AutomaticEnv()
	return c.env.GetString("redis_pass")
}

func (c *config) GetRedisDBName() int {
	c.env.AutomaticEnv()
	return c.env.GetInt("redis_db_name")
}

func NewConfig(env *viper.Viper) Config {
	return &config{
		env: env,
	}
}
