package service

import (
	"context"
	"errors"
	"webook/internal/domain"
	"webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}

func (a *articleService) update(ctx context.Context, art domain.Article) error {
	artInDB := a.repo.FindById(ctx, art.Id)
	if art.Author.Id != artInDB.Author.Id {
		return errors.New("更新别人的数据")
	}
	return a.repo.Update(ctx, art)
}
