package svc

import (
	"github.com/silenceper/wechat/v2/cache"
	"github.com/xh-polaris/auth-rpc/internal/config"
	"github.com/xh-polaris/auth-rpc/internal/model"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/miniprogram"
	mpConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type A interface {
	miniprogram.MiniProgram
}
type ServiceContext struct {
	Config      config.Config
	UserModel   model.UserModel
	Redis       *redis.Redis
	Meowchat    *miniprogram.MiniProgram
	MeowchatOld *miniprogram.MiniProgram
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(monc.MustNewModel(c.Mongo.URL, c.Mongo.DB, model.UserCollectionName, c.CacheConf)),
		Redis:     c.Redis.NewRedis(),
		Meowchat: wechat.NewWechat().GetMiniProgram(&mpConfig.Config{
			AppID:     c.Meowchat.AppID,
			AppSecret: c.Meowchat.AppSecret,
			Cache:     cache.NewMemory(),
		}),
		MeowchatOld: wechat.NewWechat().GetMiniProgram(&mpConfig.Config{
			AppID:     c.MeowchatOld.AppID,
			AppSecret: c.MeowchatOld.AppSecret,
			Cache:     cache.NewMemory(),
		}),
	}
}
