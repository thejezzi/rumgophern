package lib

import (
	"context"
	"gocontext/config"
)

func InitConfig(ctx context.Context) context.Context {
	c := config.NewSuperSecretConfig()
	c.SetValue("key1", "key2", "key", 42)
	return context.WithValue(ctx, config.MyKey, c.Unwrap())
}
