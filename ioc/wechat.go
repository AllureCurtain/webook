package ioc

import (
	"webook/internal/service/oauth2/wechat"
	"webook/internal/web"
)

func InitWechatService() wechat.Service {

	// TODO: 暂时默认写法
	var appId = "111111"
	var appKey = "1111111"

	//appId, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("WECHAT_APP_ID is not set")
	//}
	//appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("没有找到环境变量 WECHAT_APP_SECRET")
	//}
	return wechat.NewService(appId, appKey)
}

func NewWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
