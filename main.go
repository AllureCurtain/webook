package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"webook/internal/integration/startup"
)

func main() {
	//db := initDB()
	//server := initWebServer()
	//
	//rdb := initRedis()
	//u := initUser(db, rdb)
	//u.RegisterRoutes(server)

	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "你好")
	//})
	//server.Run(":8080")

	initViperV1()

	initLogger()

	server := startup.InitWebServer()

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})

	server.Run(":8080")
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	zap.L().Info("hello")
}

func initViperV1() {
	cfile := pflag.String("config", "config/config.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperV2Watch() {
	cfile := pflag.String("config",
		"config/dev.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViper() {
	viper.SetDefault("db.src.dsn", "root:root@tcp(localhost:3306)/webook")

	// 配置文件的名字，但是不包含文件扩展名
	// 不包含 .go, .yaml 之类的后缀
	viper.SetConfigName("dev")
	// 告诉 viper 我的配置用的是 yaml 格式
	viper.SetConfigType("yaml")
	// 当前工作目录下的 config 子目录
	viper.AddConfigPath("./config")

	// 读取配置到 viper 里面
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

//func initWebServer() *gin.Engine {
//	server := gin.Default()
//
//	server.Use(func(ctx *gin.Context) {
//		println("这是第一个 middleware")
//	})
//
//	server.Use(func(ctx *gin.Context) {
//		println("这是第二个 middleware")
//	})
//
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: "localhost:6379",
//	})
//
//	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
//
//	// 跨域
//	server.Use(cors.New(cors.Config{
//		//AllowOrigins: []string{"http://localhost:3000"},
//		//AllowMethods: []string{"GET", "POST"},
//		AllowHeaders: []string{"Content-Type", "Authorization"},
//		// 是否允许携带 cookie 之类的东西
//		AllowCredentials: true,
//		// 不加这个，前端拿不到
//		ExposeHeaders: []string{"x-jwt-token"},
//		AllowOriginFunc: func(origin string) bool {
//			//return origin == "https://github.com"
//			if strings.HasPrefix(origin, "http://localhost") {
//				// 开发环境
//				return true
//			}
//			return strings.Contains(origin, "company.com")
//		},
//		MaxAge: 12 * time.Hour,
//	}))
//
//	// session
//	//store := cookie.NewStore([]byte("secret"))
//
//	store := memstore.NewStore([]byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"), []byte("o6jdlo2cb9f9pb6h46fjmllw481ldebj"))
//
//	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", "",
//	//	[]byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"),
//	//	[]byte("o6jdlo2cb9f9pb6h46fjmllw481ldebj"))
//
//	//if err != nil {
//	//	panic(err)
//	//}
//
//	server.Use(sessions.Sessions("webook", store))
//
//	// 校验是否登录
//	//server.Use(middleware.NewLoginMiddleWareBuilder().
//	//	IgnorePaths("/users/signup").
//	//	IgnorePaths("/users/login").Build())
//	server.Use(middleware.NewLoginJWTMiddleWareBuilder().
//		IgnorePaths("/users/signup").
//		IgnorePaths("/users/login_sms/code/send").
//		IgnorePaths("/users/login_sms").
//		IgnorePaths("/users/login").Build())
//
//	return server
//}
