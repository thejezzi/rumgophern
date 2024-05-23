package lib

import (
	"gocontext/config"
)

func InitConfig() *config.SuperSecretConfig {
	c := config.NewSuperSecretConfig()
	c.SetValue("a", "b", "c", 1)
	return c
}
