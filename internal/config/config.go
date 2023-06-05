package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mongo struct {
		URL string
		DB  string
	}
	CacheConf cache.CacheConf
	Meowchat  struct {
		AppID     string
		AppSecret string
	}
	MeowchatOld struct {
		AppID     string
		AppSecret string
	}
	MeowchatManager struct {
		AppID     string
		AppSecret string
	}
	SMTP struct {
		Username string
		Password string
		Host     string
		Port     int
	}
}
