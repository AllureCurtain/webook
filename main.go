package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

	server := InitWebServer()

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})

	server.Run(":8080")
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
