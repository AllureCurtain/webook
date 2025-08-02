package article

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository"
	dao "webook/internal/repository/dao/article"
)

// AuthorRepository 封装user的client用于获取用户信息
type AuthorRepository interface {
	// FindAuthor id为文章id
	FindAuthor(ctx context.Context, id int64) (domain.Author, error)
}

type LocalAuthorRepository struct {
	userRepo repository.UserRepository
	dao      dao.ArticleDAO
}

// NewAuthorRepository 是 LocalAuthorRepository 的构造函数。
// 更改: 函数重命名，并且现在接收 UserService 作为参数。
func NewAuthorRepository(articleDao dao.ArticleDAO, userRepo repository.UserRepository) AuthorRepository {
	return &LocalAuthorRepository{
		userRepo: userRepo,
		dao:      articleDao,
	}
}

func (repo *LocalAuthorRepository) FindAuthor(ctx context.Context, id int64) (domain.Author, error) {
	// 首先，从数据库获取文章以找到作者的 ID。
	art, err := repo.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Author{}, err
	}

	// 更改: 不再使用 gRPC client，而是直接调用本地的 userSvc 的方法。
	// 我们假设 UserService 有一个 FindById 方法。
	u, err := repo.userRepo.FindById(ctx, art.AuthorId)
	if err != nil {
		return domain.Author{}, err
	}

	// 将 service 返回的 domain.User 转换为 domain.Author。
	return domain.Author{
		Id:   u.Id,
		Name: u.Nickname, // 假设 domain.User 中有 Nickname 字段
	}, nil
}

//type GrpcAuthorRepository struct {
//	client userv1.UserServiceClient
//	dao    dao.ArticleDAO
//}

//func NewGrpcAuthorRepository(articleDao dao.ArticleDAO, client userv1.UserServiceClient) AuthorRepository {
//	return &GrpcAuthorRepository{
//		client: client,
//		dao:    articleDao,
//	}
//}

//func (g *GrpcAuthorRepository) FindAuthor(ctx context.Context, id int64) (domain.Author, error) {
//	art, err := g.dao.GetPubById(ctx, id)
//	if err != nil {
//		return domain.Author{}, nil
//	}
//	u, err := g.client.Profile(ctx, &userv1.ProfileRequest{
//		Id: art.AuthorId,
//	})
//	if err != nil {
//		return domain.Author{}, err
//	}
//	return domain.Author{
//		Id:   u.User.Id,
//		Name: u.User.Nickname,
//	}, nil
//}
