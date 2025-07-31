package service

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/article"
	"webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid, id int64) error
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	repo article.ArticleRepository

	// V1
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository

	logger logger.LoggerV1
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func NewArticleServiceV1(author article.ArticleAuthorRepository, reader article.ArticleReaderRepository, l logger.LoggerV1) ArticleService {
	return &articleService{
		author: author,
		reader: reader,
		logger: l,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}

func (a *articleService) update(ctx context.Context, art domain.Article) error {
	return a.repo.Update(ctx, art)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, art)
}

// PublishV1 基于使用两种 repository 的写法
func (svc *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	// 这一段逻辑其实就是 Save
	if art.Id == 0 {
		id, err = svc.author.Create(ctx, art)
	} else {
		err = svc.author.Update(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	// 保持制作库和线上库的 ID 是一样的。
	art.Id = id
	for i := 0; i < 3; i++ {
		err = svc.reader.Save(ctx, art)
		if err == nil {
			break
		}
		svc.logger.Error("部分失败：保存数据到线上库失败",
			logger.Field{Key: "art_id", Value: id},
			logger.Error(err))
		// 在接入了 metrics 或者 tracing 之后，
		// 这边要进一步记录必要的DEBUG信息。
	}
	if err != nil {
		//svc.logger.Error("部分失败：保存数据到线上库重试都失败了",
		//	logger.Field{Key: "art_id", Value: id},
		//	logger.Error(err))
		return 0, err
	}
	return id, nil
}

func (svc *articleService) Withdraw(ctx context.Context, uid, id int64) error {
	return svc.repo.SyncStatus(ctx, uid, id, domain.ArticleStatusPrivate)
}
