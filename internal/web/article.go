package web

import (
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
	"webook/internal/domain"
	"webook/internal/service"
	//ijwt "webook/internal/web/jwt"
	"webook/pkg/ginx"
	"webook/pkg/logger"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc     service.ArticleService
	l       logger.LoggerV1
	intrSvc service.InteractiveService
	reward  service.RewardService
	biz     string
}

func NewArticleHandler(svc service.ArticleService, l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
		biz: "article",
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", h.Edit)
	g.POST("/withdraw", h.Withdraw)
	g.POST("/publish", h.Publish)

	g.POST("/list", ginx.WrapClaimsAndReq(h.List))
	g.GET("/detail/:id", h.Detail)

	pub := g.Group("/pub")
	//pub.GET("/pub", a.PubList)
	pub.GET("/:id", ginx.WrapClaims(h.PubDetail))
	pub.POST("/like", ginx.WrapClaimsAndReq[LikeReq](h.Like))
	pub.POST("/collect", ginx.WrapClaimsAndReq[CollectReq](h.Collect))
	// 打赏
	//pub.POST("/reward", ginx.WrapClaimsAndReq[RewardReq](h.Reward))
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	c := ctx.MustGet("claims")
	claims, ok := c.(*ginx.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("未发现用户的 session 信息")
		return
	}

	// 检测输入，跳过这一步
	// 调用 svc 的代码
	id, err := h.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.Id,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 日志
		h.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ginx.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("未发现用户的 session 信息")
		return
	}
	// 检测输入，跳过这一步
	// 调用 svc 的代码
	id, err := h.svc.Publish(ctx, req.toDomain(claims.Id))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 日志
		h.l.Error("发表帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ginx.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("未发现用户的 session 信息")
		return
	}
	// 检测输入，跳过这一步
	// 调用 svc 的代码
	err := h.svc.Withdraw(ctx, claims.Id, req.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 日志
		h.l.Error("撤回帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: req.Id,
	})
}

type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (a *ArticleHandler) List(ctx *gin.Context, req ListReq, uc ginx.UserClaims) (ginx.Result, error) {

	// 对于批量接口来说，要小心批次大小
	//if req.Limit > 100 {
	//	a.l.Error("获得用户会话信息失败，LIMIT过大")
	//	return ginx.Result{
	//		Code: 4,
	//		// 我会倾向于不告诉攻击者批次太大
	//		// 因为一般你和前端一起完成任务的时候
	//		// 你们是协商好了的，所以会进来这个分支
	//		// 就表明是有人跟你过不去
	//		Msg: "请求有误",
	//	}, nil
	//}

	res, err := a.svc.List(ctx, uc.Id, req.Offset, req.Limit)
	if err != nil {
		a.l.Error("获得用户会话信息失败")
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, nil
	}
	return ginx.Result{
		Data: slice.Map[domain.Article, ArticleVo](res,
			func(idx int, src domain.Article) ArticleVo {
				return ArticleVo{
					Id:       src.Id,
					Title:    src.Title,
					Abstract: src.Abstract(),
					Status:   src.Status.ToUint8(),
					// 这个列表请求，不需要返回内容
					//Content: src.Content,
					// 这个是创作者看自己的文章列表，也不需要这个字段
					//Author: src.Author
					Ctime: src.Ctime.Format(time.DateTime),
					Utime: src.Utime.Format(time.DateTime),
				}
			}),
	}, nil
}

func (a *ArticleHandler) Detail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		a.l.Error("前端输入的 ID 不对", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("user").(ginx.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	art, err := a.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得文章信息失败", logger.Error(err))
		return
	}
	//art := resp.GetArticle()
	// 这是不借助数据库查询来判定的方法
	if art.Author.Id != usr.Id {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			// 也不需要告诉前端究竟发生了什么
			Msg: "输入有误",
		})
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。
		a.l.Error("非法访问文章，创作者 ID 不匹配", logger.Int64("uid", usr.Id))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVo{
			Id:    art.Id,
			Title: art.Title,
			// 不需要这个摘要信息
			//Abstract: art.Abstract(),
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			// 这个是创作者看自己的文章列表，也不需要这个字段
			//Author: art.Author
			Ctime: art.Ctime.Format(time.DateTime),
			Utime: art.Utime.Format(time.DateTime),
		},
	})
}

func (a *ArticleHandler) PubDetail(ctx *gin.Context, ucc ginx.UserClaims) (Result, error) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		a.l.Error("前端输入的 ID 不对", logger.Error(err))
		return Result{
			Code: 4,
			Msg:  "参数错误",
		}, fmt.Errorf("查询文章详情的 ID %s 不正确, %w", idstr, err)
	}

	var eg errgroup.Group
	var art domain.Article
	eg.Go(func() error {
		art, err = a.svc.GetPublishedById(ctx, id)
		return err
	})

	var intr domain.Interactive
	// 获得这篇文章的计数
	eg.Go(func() error {
		uc := ctx.MustGet("users").(ginx.UserClaims)
		intr, err = a.intrSvc.Get(ctx, a.biz, id, uc.Id)
		return err
	})

	err = eg.Wait()
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, fmt.Errorf("查询出错")
	}

	// 直接异步操作，在确定我们获取到了数据之后再来操作
	go func() {
		err = a.intrSvc.IncrReadCnt(ctx, a.biz, art.Id)
		if err != nil {
			a.l.Error("增加文章阅读数失败", logger.Error(err))
		}
	}()

	return Result{
		Data: ArticleVo{
			Id:      art.Id,
			Title:   art.Title,
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			// 要把作者信息带出去
			Author:     art.Author.Name,
			Ctime:      art.Ctime.Format(time.DateTime),
			Utime:      art.Utime.Format(time.DateTime),
			ReadCnt:    intr.ReadCnt,
			CollectCnt: intr.CollectCnt,
			LikeCnt:    intr.LikeCnt,
			Liked:      intr.Liked,
			Collected:  intr.Collected,
		},
	}, nil
}

func (a *ArticleHandler) Like(ctx *gin.Context, req LikeReq, uc ginx.UserClaims) (ginx.Result, error) {
	var err error
	if req.Like {
		err = a.intrSvc.Like(ctx, a.biz, req.Id, uc.Id)
	} else {
		err = a.intrSvc.CancelLike(ctx, a.biz, req.Id, uc.Id)
	}

	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{Msg: "OK"}, nil
}

func (a *ArticleHandler) Collect(ctx *gin.Context, req CollectReq, uc ginx.UserClaims) (Result, error) {
	err := a.intrSvc.Collect(ctx, a.biz, req.Id, uc.Id, req.Cid)
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{Msg: "OK"}, nil
}

//func (a *ArticleHandler) Reward(ctx *gin.Context, req RewardReq, uc ginx.UserClaims) (ginx.Result, error) {
//	art, err := a.svc.GetPublishedById(ctx.Request.Context(), req.Id)
//	if err != nil {
//		return ginx.Result{
//			Code: 5,
//			Msg:  "系统错误",
//		}, err
//	}
//	resp, err := a.reward.PreReward(ctx.Request.Context(), &rewardv1.PreRewardRequest{
//		Biz:       "article",
//		BizId:     art.Id,
//		BizName:   art.Title,
//		TargetUid: art.Author.GetId(),
//		Uid:       uc.Id,
//		Amt:       req.Amt,
//	})
//	if err != nil {
//		return ginx.Result{
//			Code: 5,
//			Msg:  "系统错误",
//		}, err
//	}
//	return ginx.Result{
//		Data: map[string]any{
//			"codeURL": resp.CodeUrl,
//			"rid":     resp.Rid,
//		},
//	}, nil
//}

func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}
